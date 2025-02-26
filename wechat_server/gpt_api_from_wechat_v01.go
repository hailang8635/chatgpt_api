package wechat_server

import (
	"chatgpt_api/api_from_ai"
	"chatgpt_api/domain"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

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
						fmt.Fprintf(w, "%s", MakeResponseString(toUserName, fromUserName, keywordsInfo.Answer))
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
							fmt.Fprintf(w, "%s", MakeResponseString(toUserName, fromUserName, notReturnedKeyword.Answer))
						} else {
							// 否则返回空
							fmt.Fprintf(w, "%s", MakeResponseString(toUserName, fromUserName, ""))
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
				respStr, err := api_from_ai.GptApi(keywordString)
				if err != nil {
					fmt.Fprintf(w, "%s", MakeResponseString(toUserName, fromUserName, "系统忙，请稍后再试."))
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

				fmt.Fprintf(w, "%s", MakeResponseString(toUserName, fromUserName, respStr))

			} else {
				// 浏览器直接访问的
				chatGptProcess(w, r)
			}
		}

	})
}
