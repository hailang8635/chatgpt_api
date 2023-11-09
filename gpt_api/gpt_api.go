package gpt_api

import (
    "bufio"
    "encoding/xml"
    "fmt"
    "chatgpt_api/config"
    "chatgpt_api/domain"
    "chatgpt_api/utils"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
    "time"
)

//user_word := make(map[string] string)

var timeLayoutStr = "2006-01-02 15:04:05"

var keywords = map[string]domain.RespMsg{}

func Gpt_http_server() {

    // URL请求方式
    chatGptHandler()

    // 输入输出落库版本
    chatHandlerWithDB()

    //chatHandler_bak()

    log.Println("Starting server...")
    http.ListenAndServe(":8080", nil)
}

/*
*

    GPT3.5
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

        fmt.Fprintf(w, "%s", "请输入您要问的内容？")
        return
    }

    if msgType != "text" {
        log.Println("非text消息", toUserName, fromUserName, msgType, voiceText)
        if msgType == "voice" && voiceText != "" {
            keywordParamsOrigin = voiceText
        } else {
            fmt.Fprintf(w, "%s", makeResponseString(toUserName, fromUserName, "目前我只能回答文字内容.."))
            return
        }
    }

    // 过滤敏感关键词
    keywordParams := utils.Substring(keywordParamsOrigin, 20)

    if config.VerfiyBadWordsOnlyResult(keywordParamsOrigin) {
        fmt.Fprintf(w, "%s", makeResponseString(toUserName, fromUserName, "该问题受限于法律法规限制无法回答.."))
        return
    }

    // 微信限制5s返回，5s未返回有3次重试
    if strings.Contains("1?？。，.,", keywordParamsOrigin) {
        _, keywordItems := utils.SelectOne(domain.Keywords{
            Fromuser: fromUserName,
            //Is_done: 1,
            OrderByIdDesc: true,
        })
        keywordParamsOrigin = keywordItems.Keyword
        keywordParams = keywordItems.Keyword
    }

    // 取mysql中数据
    //keywordsInfo, exists := keywords[keywordParamsOrigin]
    rows, keywordInDb := utils.SelectOne(domain.Keywords{
        Keyword: keywordParamsOrigin,
    })

    // A=已存在的关键字
    if rows >= 1 {
        processExistsKeyword(w, keywordInDb, keywordParams, fromUserName, toUserName)
    } else {
        processNewKeyword(w, keywordParamsOrigin, keywordParams, fromUserName, toUserName, startTime)
    }
}

/**
 * B流程，新的请求
 */
func processNewKeyword(w http.ResponseWriter, keywordParamsOrigin string, keywordParams string, fromUserName string, toUserName string, startTime time.Time) {
    // B = 第一次查询的关键字则放进map

    // B1 = 第一次查询，关键字先入库
    lastId := utils.Insert(domain.Keywords{
        Fromuser:    fromUserName,
        Keyword:     keywordParamsOrigin,
        Answer:      "",
        Labels:      "",
        Catalog:     "",
        Is_done:     2,
        Is_finished: 2,
        Create_time: startTime,
        Finish_time: startTime,
    })

    offset_5m, _ := time.ParseDuration("-1m")
    _, userHistoryMessage := utils.SelectList(domain.Keywords{
        Fromuser:          fromUserName,
        Create_time_start: time.Now().Add(offset_5m),
        //Is_done:     1,
        //Is_finished: 2,
        //Keyword:  keywordParamsOrigin,
    }, 3)

    // B2 = 发起调用gpt的api
    // 根据关键词查询GPT接口
    apiStart := time.Now()
    log.Printf("B2 开始查询openai.com %s \n", keywordParams)
    respStr, err := GptApi2(keywordParamsOrigin, userHistoryMessage)

    if err != nil {
        fmt.Fprintf(w, "%s", makeResponseString(toUserName, fromUserName, "系统忙，请稍后再试."))
        return
    }

    // 微信最大2048字节
    respStr = utils.SubstringByBytes(respStr, 2040)
    isBad, respStrModified := config.VerfiyBadWords(respStr)
    if isBad {
        // fmt.Fprintf(w, "%s", makeResponseString(toUserName, fromUserName, "该问题受限于法律法规限制无法回答.."))
        respStr = respStrModified
    }

    // 保存记录，超过15s的为未返回状态，小于15s的为已返回状态
    endTime := time.Now()

    // B3 = 调用gpt api结束
    log.Printf("B3 查询openai.com成功 %s 耗时 %d s \n", keywordParams, endTime.Unix()-apiStart.Unix())

    timeSpend := endTime.Unix() - startTime.Unix()
    is_finished := 1

    // 1 * retry_gap ?
    if timeSpend < 1*retry_gap {
        fmt.Fprintf(w, "%s", makeResponseString(toUserName, fromUserName, respStr))
    } else {
        is_finished = 2
    }

    keywords_final := domain.Keywords{
        Id:          lastId,
        Fromuser:    fromUserName,
        Keyword:     keywordParamsOrigin,
        Labels:      "",
        Catalog:     "",
        Create_time: startTime,
        Answer:      respStr,
        Is_done:     1,
        Is_finished: is_finished,
        Finish_time: endTime,
    }

    // B4 查询openai.com成功更新答案及is_done状态
    utils.Update(keywords_final)
    log.Printf("<---- B4 更新状态结束 keywordParams: %s, is_done: %d, is_finished: %d, 流程耗时: %d s \n\n", keywordParams, 1, is_finished, timeSpend)
    return
}

var retry_gap int64 = 5

/**
 * A流程，wx轮询
 */
func processExistsKeyword(w http.ResponseWriter, keywordInDb domain.Keywords, keywordParams string, fromUserName string, toUserName string) {
    // A1 = 已完成
    if keywordInDb.Is_done == 1 {
        log.Printf("<---- A1 直接返回已完成的keyword： %s", keywordParams)
        fmt.Fprintf(w, "%s", makeResponseString(toUserName, fromUserName, keywordInDb.Answer))

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
        _, keywordInDbAt15s := utils.SelectOne(domain.Keywords{
            Keyword: keywordInDb.Keyword,
        })

        if keywordInDbAt15s.Is_done == 1 {
            // 存在之前已完成未返回的记录
            log.Printf("<---- A2.1 wechat retry 3 ... >12s的请求(%d s) 该用户有已查得未返回的keyword %s \n", time_spend, keywordInDbAt15s.Keyword)

            // 返回未完成的记录，并更新记录的is_finished状态
            fmt.Fprintf(w, "%s", makeResponseString(toUserName, fromUserName, keywordInDbAt15s.Answer))

            keywordInDbAt15s.Is_finished = 1
            utils.Update(keywordInDbAt15s)
            return
        }

        // A2.2 = 临界15s时渠道仍未返回
        // 查找该用户已完成且未返回的记录
        not_returned_rows, keywordInDb_not_returned := utils.SelectOne(domain.Keywords{
            Fromuser:    fromUserName,
            Is_done:     1,
            Is_finished: 2,
            //Keyword:  keywordParamsOrigin,
        })

        if not_returned_rows >= 1 {
            // 存在之前已完成未返回的记录
            log.Printf("<---- A2.2 wechat retry 3 ... >12s的请求(%d s) 该用户有已查得未返回的keyword %s \n", time_spend, keywordInDb_not_returned.Keyword)

            // 返回未完成的记录，并更新记录的is_finished状态
            fmt.Fprintf(w, "%s", makeResponseString(toUserName, fromUserName, keywordInDb_not_returned.Answer))

            keywordInDb_not_returned.Is_finished = 1
            utils.Update(keywordInDb_not_returned)
            return
        } else {
            // 15s内未查成功，且无未返回的记录时

            log.Printf("<---- A2.3 关键字正在处理中(已耗时:%d ), 回复给client进行重试 %s \n", time_spend, keywordParams)

            // 收到粉丝消息后不想或者不能5秒内回复时，需回复“success”字符串（下文详细介绍）
            //fmt.Fprintf(w, "%s", makeResponseString(toUserName, fromUserName, "结果生成中...，请5s后再问一遍"))
            fmt.Fprintf(w, "%s", makeResponseString(toUserName, fromUserName, "答案生成中, 请5s后回复【1】获取答案"))
        }
        // return
    }
    return
}

func chatGptHandler() {
    http.HandleFunc("/chatgpt_20230401", func(w http.ResponseWriter, r *http.Request) {
        chatGptProcess(w, r)
    })
}

func chatGptProcess(w http.ResponseWriter, r *http.Request) {
    s := r.URL.Query().Get("s")
    if s == "" || len(s) <= 0 {
        // 关键词为空时返回
        fmt.Fprintf(w, "%s", "请输入您要问的内容？")
        return
    } else {
        // 根据关键词查询GPT接口
        respStr, err := GptApi(s)
        if err != nil {
            fmt.Fprintf(w, "%s %s", s, "系统忙，请稍后再试.")
            return
        }
        fmt.Fprintf(w, "【问】：%s \n\n【答】: \n %s", s, respStr)
    }
}

func chatHandler_bak() {
    http.HandleFunc("/chat_20230326", func(w http.ResponseWriter, r *http.Request) {
        // 校验微信平台-验证开发平台
        echostr := r.URL.Query().Get("echostr")
        if echostr != "" && len(echostr) >= 1 {
            fmt.Fprintf(w, "%s", echostr)
            return
        }

        // 处理微信订阅号聊天窗口中的问题【xml】
        // 用户 --> 微信平台 --> okzhang.com服务端 --> gpt
        data, err := io.ReadAll(r.Body)
        //log.Println("data <----", data)

        if err != nil {
            log.Println("io.ReadAll error")
        } else if data != nil {
            // 请求xml
            reqInfo := domain.WXReqTextMsg{}
            err := xml.Unmarshal(data, &reqInfo)
            if err != nil {
                log.Println("xml.Unmarshal error")
            }

            // 请求参数
            fromUserName := reqInfo.FromUserName
            toUserName := reqInfo.ToUserName
            keywordString := reqInfo.Content
            msgType := reqInfo.MsgType
            msgId := reqInfo.MsgId
            createTime := reqInfo.CreateTime

            log.Printf("reqInfo: %s, from: %s, msgId: %s, msgType: %s, time: %s, keywordString: %s \n", reqInfo, fromUserName, msgId, msgType, createTime, keywordString)

            // TODO 微信限制5s返回，5s未返回有3次重拾
            // 相同关键词不再查询
            if keywordString != "" {
                keywordsInfo, exists := keywords[keywordString]
                if exists {
                    if keywordsInfo.IsDone == 1 {
                        log.Println("存在的keyword，且已完成，直接返回", keywordString)
                        fmt.Fprintf(w, "%s", makeResponseString(toUserName, fromUserName, keywordsInfo.Answer))
                    } else {
                        log.Println("存在的keyword，未完成，返回空", keywordString)
                        // TODO 查找该用户已完成且未返回的记录
                        existNotReturnedAnswer := false
                        var notReturnedKeyword domain.RespMsg
                        for _, v := range keywords {
                            if v.Username == fromUserName && v.IsDone == 1 && v.IsReturned != 1 {
                                existNotReturnedAnswer = true
                                notReturnedKeyword = v
                                v.IsReturned = 1
                            }
                        }

                        if existNotReturnedAnswer {
                            fmt.Fprintf(w, "%s", makeResponseString(toUserName, fromUserName, notReturnedKeyword.Answer))
                        } else {
                            // 否则返回空
                            fmt.Fprintf(w, "%s", makeResponseString(toUserName, fromUserName, ""))
                        }
                    }
                    return
                } else {
                    // 第一次查询则放进map
                    keywords[keywordString] = domain.RespMsg{
                        Username:   fromUserName,
                        Keyword:    keywordString,
                        IsDone:     2,
                        IsReturned: 2,
                    }
                }

                // 根据关键词查询GPT接口
                // 有该用户未返回消息则返回
                respStr, err := GptApi(keywordString)
                if err != nil {
                    fmt.Fprintf(w, "%s", makeResponseString(toUserName, fromUserName, "系统忙，请稍后再试."))
                    return
                }

                // 返回给微信结果
                keywords[keywordString] = domain.RespMsg{
                    Username:   fromUserName,
                    Keyword:    keywordString,
                    Answer:     respStr,
                    IsDone:     1,
                    IsReturned: 2,
                }

                fmt.Fprintf(w, "%s", makeResponseString(toUserName, fromUserName, respStr))

            } else {
                // 浏览器直接访问的
                chatGptProcess(w, r)
            }
        }

    })
}

func makeResponseString(toUserName string, fromUserName string, respStr string) string {
    return makeResponseString2(toUserName, fromUserName, "text", respStr)
}
func makeResponseString2(toUserName string, fromUserName string, msgType string, respStr string) string {
    respInfo := domain.WXRespTextMsg{}
    respInfo.FromUserName = domain.CDATA{toUserName}
    respInfo.ToUserName = domain.CDATA{fromUserName}
    respInfo.MsgType = domain.CDATA{msgType}
    respInfo.Content = domain.CDATA{respStr}
    respInfo.CreateTime = time.Now().Unix()

    respXml2String, _ := xml.MarshalIndent(respInfo, "", "")
    return string(respXml2String)
}

/*
命令行使用
*/
func Ask_gpt() (string, error) {
    if len(os.Args) >= 2 {
        args1 := os.Args[1]
        respBody, err := GptApi(args1)

        if err != nil {
            log.Println("gpt_api 调用失败")
            return "", err
        }
        log.Println(respBody)
        log.Println("请输入您要问的内容？【N】退出")

    } else {
        log.Println("请输入您要问的内容？")
    }
    input := bufio.NewScanner(os.Stdin)
    for input.Scan() {
        line := input.Text()

        // 遇到N退出
        if line == "N" {
            break
        }

        log.Println("正在问ChatGpt ... 【", line, "】")
        resp, err := GptApi(line)
        if err != nil {
            log.Println("gpt_api 调用失败")
            return "", err
        } else {
            log.Println(resp)
        }
        //log.Println()
        fmt.Println()
        log.Println("请继续输入您要问的内容？")
        fmt.Println()
    }

    return "ask done", nil
}
