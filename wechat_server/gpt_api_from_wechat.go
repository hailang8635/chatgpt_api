package wechat_server

import (
	"chatgpt_api/config"
	"chatgpt_api/domain"
	"chatgpt_api/utils"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//user_word := make(map[string] string)

var timeLayoutStr = "2006-01-02 15:04:05"

//var timeLayoutStrYYYYMMDDHHmmss = "20060102150405"

//var keywords = map[string]domain.RespMsg{}
var retry_gap int64 = 5
var length_wechat = 500

//var length_wechat = 300

// TODO 入库等操作独立出来
func Gpt_http_server() {

	// URL请求方式
	chatGptHandler()

	// 输入输出落库版本
	chatHandlerWithDB()

	//chatHandler_bak()

	log.Println("Starting server...")
	// ":8080"
	http.ListenAndServe(config.DefaultPort, nil)
}

/**
 *
 * 微信后台发来的请求
 */
func chatHandlerWithDB() {
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Println("")

		// 校验微信平台-验证开发平台
		echostr := r.URL.Query().Get("echostr")
		if echostr != "" && len(echostr) >= 1 {
			fmt.Fprintf(w, "%s", echostr)
			return
		}

		// TODO 支持GPT-4
		openAiModelVersion := r.URL.Query().Get("gpt_api_version")
		log.Println("openAiModelVersion", openAiModelVersion)

		//
		// 用户 --> 微信平台 --> okzhang.com服务端 --> gpt
		data, err := io.ReadAll(r.Body)
		//log.Println("data <----", data)

		if err != nil {
			log.Println("io.ReadAll error")
			return
		}
		if data != nil {
			processWechatRequest(w, r, data, startTime)
		}

	})
}

/**
 * 处理微信订阅号聊天窗口中的请求【xml】
 */
func processWechatRequest(w http.ResponseWriter, r *http.Request, data []byte, startTime time.Time) {
	// 请求xml
	reqInfo := domain.WXReqTextMsg{}
	err := xml.Unmarshal(data, &reqInfo)
	if err != nil {
		log.Println("xml.Unmarshal error", err)
	}

	// 请求参数
	fromUserName := reqInfo.FromUserName
	toUserName := reqInfo.ToUserName
	keywordParamsOrigin := reqInfo.Content
	msgId := reqInfo.MsgId
	createTime := reqInfo.CreateTime
	msgType := reqInfo.MsgType
	voiceText := reqInfo.Recognition

	log.Printf("----> A0 toUserName: %s, from: %s, msgId: %d, msgType: %s time: %s, keywordParams: %s ",
		toUserName, fromUserName, msgId, msgType, time.Unix(createTime, 0).Format(timeLayoutStr), keywordParamsOrigin)

	if keywordParamsOrigin == "" && msgType == "" {
		// 浏览器直接访问的
		// chatGptProcess(w, r)

		fmt.Fprintf(w, "%s", "请输入您要问的内容 ？")
		return
	}

	if msgType != "text" {
		log.Println("非text消息", toUserName, fromUserName, msgType, voiceText)
		if msgType == "voice" && voiceText != "" {
			keywordParamsOrigin = voiceText
		} else {
			fmt.Fprintf(w, "%s", MakeResponseString(toUserName, fromUserName, "目前我只能回答文字内容.."))
			return
		}
	}

	// 输入参数 过滤敏感关键词
	keywordParams := utils.Substring(keywordParamsOrigin, 20)

	if config.VerfiyBadWordsOnlyResult(keywordParamsOrigin) {
		fmt.Fprintf(w, "%s", MakeResponseString(toUserName, fromUserName, "该问题受限于法律法规限制无法回答.."))
		return
	}

	// 微信限制5s返回，5s未返回有3次重试
	if strings.Contains("1?？。，.,", keywordParamsOrigin) {
		_, keywordItems := utils.SelectOne(domain.KeywordAndAnswerItem{
			Fromuser: fromUserName,
			//Is_done: 1,
			OrderByIdDesc: true,
		})
		keywordParamsOrigin = keywordItems.Keyword
		keywordParams = keywordItems.Keyword
	}

	// 取mysql中数据
	offset_5m, _ := time.ParseDuration("-24h")
	//keywordsInfo, exists := keywords[keywordParamsOrigin]
	rows, keywordAnAnswerInDb := utils.SelectOne(domain.KeywordAndAnswerItem{
		Keyword:           keywordParamsOrigin,
		Create_time_start: time.Now().Add(offset_5m),
	})

	// A=已存在的关键字
	if rows >= 1 {
		processExistsKeyword(w, keywordAnAnswerInDb, keywordParams, fromUserName, toUserName)
	} else {
		processNewKeyword(w, keywordParamsOrigin, keywordParams, fromUserName, toUserName, startTime)
	}
}

/**
 * 新的请求(命名为B流程-新请求)

 * keywordParamsOrigin 原始关键字
 * keywordParams 排除敏感词的关键字
 */
func processNewKeyword(w http.ResponseWriter, keywordParamsOrigin string, keywordParams string, fromUserName string, toUserName string, startTime time.Time) {
	// B = 第一次查询的关键字则放进map

	// B1 = 第一次查询，关键字先入库
	lastId, userHistoryMessage := InsertItemAndReturnHistory(fromUserName, keywordParamsOrigin, startTime)

	// B2 = 发起调用gpt的api
	// 根据关键词查询GPT接口
	apiStart := time.Now()
	log.Printf("B2 开始查询openai.com %s \n", keywordParams)
	respStr, err := GetAPIResult(keywordParamsOrigin, userHistoryMessage)

	if err != nil {
		fmt.Fprintf(w, "%s", MakeResponseString(toUserName, fromUserName, "系统忙，请稍后再试."))
		return
	}

	longStringUrl := SaveAsHTML(respStr, keywordParamsOrigin, startTime)

	// 微信最大2048字节
	//respStr = utils.SubstringByBytes(respStr, 2000-len(longStringUrl)) + "\n" + longStringUrl

	isBad, respStrModified := config.VerfiyBadWords(respStr)
	if isBad {
		// fmt.Fprintf(w, "%s", makeResponseString(toUserName, fromUserName, "该问题受限于法律法规限制无法回答.."))
		respStr = respStrModified
	}

	responseString := respStr
	if len(respStr) > length_wechat {
		responseString = "[答案详情见链接] \n" + config.HtmlUrl + longStringUrl
	}

	// 保存记录，超过15s的为未返回状态，小于15s的为已返回状态
	endTime := time.Now()

	// B3 = 调用gpt api结束
	log.Printf("B3 查询ai接口成功 %s 耗时 %d s \n", keywordParams, endTime.Unix()-apiStart.Unix())

	timeSpend := endTime.Unix() - startTime.Unix()
	is_finished := 1

	// 1 * retry_gap ?
	if timeSpend < 1*retry_gap {
		fmt.Fprintf(w, "%s", MakeResponseString(toUserName, fromUserName, responseString))
	} else {
		is_finished = 2
	}

	UpdateItem(lastId, fromUserName, keywordParamsOrigin, startTime, respStr, longStringUrl, is_finished, endTime)
	log.Printf("<---- B4 更新状态结束 keywordParams: %s, is_done: %d, is_finished: %d, 流程耗时: %d s \n\n", keywordParams, 1, is_finished, timeSpend)
	return
}

/**
 * wx轮询&查询老的问题(命名为A流程-老请求)
 */
func processExistsKeyword(w http.ResponseWriter, keywordInDb domain.KeywordAndAnswerItem, keywordParams string, fromUserName string, toUserName string) {

	urlString := ""
	if keywordInDb.Url != "" {
		urlString = "[答案详情见链接]\n" + config.HtmlUrl + url.QueryEscape(keywordInDb.Url)
	}

	// A1 = 已完成
	if keywordInDb.Is_done == 1 {
		log.Printf("<---- A1 直接返回已完成的keyword： %s", keywordParams)

		is_repeat_question := ""
		if keywordInDb.Is_finished == 1 {
			is_repeat_question = "\n[重复问题]"
		}
		fmt.Fprintf(w, "%s", MakeResponseString(toUserName, fromUserName, keywordInDb.Answer+"\n"+urlString+is_repeat_question))

		// 对应更新为已返回
		if keywordInDb.Is_finished != 1 {
			keywordInDb.Is_finished = 1
			utils.Update(keywordInDb)
		}
		// return
	} else {
		// A2 = 渠道尚未返回结果

		// A2.0 = WX第2次查询
		// 前2次轮询不做应答，第3次返回请重试
		// 10s内不对WX查询做应答，10s后的轮询返回请重试
		time_spend := time.Now().Unix() - keywordInDb.Create_time.Unix()
		if time_spend < 2*retry_gap {
			// 第2次查询
			log.Printf("<---- A2.0 wechat retry 2... <10s的请求(%d s) 关键字正在处理中(status_code:504) %s \n", time_spend, keywordParams)

			time.Sleep(time.Duration(retry_gap) * time.Second)

			w.WriteHeader(504)
			fmt.Fprintf(w, "%s", "success")

			// 阻止走A2.1的流程
			return
		}

		// A2.1 = WX 第3次查询（超过10s的)直接拖到14s再返回
		time.Sleep(time.Duration(float32(retry_gap)-1.5) * time.Second)

		// 14s时再查一次结果
		_, keywordInDbAt15s := utils.SelectOne(domain.KeywordAndAnswerItem{
			Keyword: keywordInDb.Keyword,
		})

		if keywordInDbAt15s.Is_done == 1 {
			// 存在之前已完成未返回的记录
			log.Printf("<---- A2.1 wechat retry 3 ... >12s的请求(%d s) 该用户有已查得未返回的keyword %s \n", time_spend, keywordInDbAt15s.Keyword)

			// 返回未完成的记录，并更新记录的is_finished状态
			fmt.Fprintf(w, "%s", MakeResponseString(toUserName, fromUserName, keywordInDbAt15s.Answer+"\n"+urlString))

			keywordInDbAt15s.Is_finished = 1
			utils.Update(keywordInDbAt15s)
			return
		}

		// A2.2 = 临界15s时渠道仍未返回
		// 查找该用户已完成且未返回的记录
		not_returned_rows, keywordInDb_not_returned := utils.SelectOne(domain.KeywordAndAnswerItem{
			Fromuser:    fromUserName,
			Is_done:     1,
			Is_finished: 2,
			//Keyword:  keywordParamsOrigin,
		})

		if not_returned_rows >= 1 {
			// 存在之前已完成未返回的记录
			log.Printf("<---- A2.2 wechat retry 3 ... >12s的请求(%d s) 该用户有已查得未返回的keyword %s \n", time_spend, keywordInDb_not_returned.Keyword)

			// 返回未完成的记录，并更新记录的is_finished状态
			fmt.Fprintf(w, "%s", MakeResponseString(toUserName, fromUserName, keywordInDb_not_returned.Answer+"\n"+urlString))

			keywordInDb_not_returned.Is_finished = 1
			utils.Update(keywordInDb_not_returned)
			return
		} else {
			// 15s内未查成功，且无未返回的记录时

			log.Printf("<---- A2.3 关键字正在处理中(已耗时:%d ), 回复给client进行重试 %s \n", time_spend, keywordParams)

			// 收到粉丝消息后不想或者不能5秒内回复时，需回复“success”字符串（下文详细介绍）
			fmt.Fprintf(w, "%s", MakeResponseString(toUserName, fromUserName, "答案生成中, 请15s后回复【1】获取答案"))
		}
		// return
	}
	return
}
