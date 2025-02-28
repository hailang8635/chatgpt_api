package wechat_server

import (
	"bytes"
	"chatgpt_api/config"
	"chatgpt_api/domain"
	"chatgpt_api/utils"
	"encoding/xml"
	"github.com/yuin/goldmark"
	"log"
	"os"
	"strings"
	"time"
)

var timeLayoutStrYYYYMMDDHHmmss = "20060102150405"

var keywords = map[string]domain.RespMsg{}

var max_length_wechat = 2000

func InsertItemAndReturnHistory(fromUserName string, keywordParamsOrigin string, startTime time.Time) (int64, []domain.KeywordAndAnswerItem) {
	lastId := utils.Insert(domain.KeywordAndAnswerItem{
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

	// 查询该用户的5分钟内的历史记录
	offset_5m, _ := time.ParseDuration("-1m")
	_, userHistoryMessage := utils.SelectList(domain.KeywordAndAnswerItem{
		Fromuser:          fromUserName,
		Create_time_start: time.Now().Add(offset_5m),
		//Is_done:     1,
		//Is_finished: 2,
		//Keyword:  keywordParamsOrigin,
	}, 3)
	return lastId, userHistoryMessage
}

func UpdateItem(lastId int64, fromUserName string, keywordParamsOrigin string, startTime time.Time, respStr string, longStringUrl string, is_finished int, endTime time.Time) {
	keywords_final := domain.KeywordAndAnswerItem{
		Id:          lastId,
		Fromuser:    fromUserName,
		Keyword:     keywordParamsOrigin,
		Labels:      "",
		Catalog:     "",
		Create_time: startTime,
		Answer:      respStr,
		Url:         longStringUrl,
		Is_done:     1,
		Is_finished: is_finished,
		Finish_time: endTime,
	}

	// B4 查询openai.com成功更新答案及is_done状态
	utils.Update(keywords_final)
}

func SaveAsHTML(respStr string, keywordParamsOrigin string, startTime time.Time) string {
	longStringUrl := ""
	if len(respStr) > max_length_wechat || strings.Contains(respStr, "```") {
		var buf bytes.Buffer
		err := goldmark.Convert([]byte(respStr), &buf)
		if err != nil {
			log.Println("markdown --> html, exception", err)
		} else {
			log.Println("markdown --> html, ", utils.Substring(buf.String(), 20))

			fileNameRight := utils.Substring(strings.ReplaceAll(strings.ReplaceAll(keywordParamsOrigin, " ", ""), "/", ""), 12)
			htmlFile := startTime.Format(timeLayoutStrYYYYMMDDHHmmss) + "_" + fileNameRight + ".html"
			//htmlUrlPath := startTime.Format(timeLayoutStrYYYYMMDDHHmmss) + "_" + url.QueryEscape(fileNameRight) + ".html"

			file, err := os.Create(config.HtmlDir + htmlFile)
			if err != nil {
				log.Println("create html file error", err)
			}
			file.WriteString("<html><head>  " +
				"<title>ChatGPT助手-安德鲁家的550W</title>  " +
				"<basefont face=\"微软雅黑\" size=\"2\" />  " +
				"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no\">" +
				"<meta http-equiv=\"Content-Type\" content=\"text/html;charset=utf-8\" />" +
				"</head>")

			file.WriteString("<h4>有用户提问如下：</h4>\n")

			file.WriteString("<h2>")
			file.WriteString(keywordParamsOrigin)
			file.WriteString("</h2>\n\n\n")

			file.WriteString("<h4>以下是来自" + strings.ToUpper(config.DefaultAPI) + "的回答：</h4>\n")

			defer file.Close()

			_, err = buf.WriteTo(file)
			if err != nil {
				log.Println("write html file error", err)
			}

			// https://chatapi.okzhang.com/html/cah/test.html
			// "[答案详情见链接] \n" + utils.HtmlUrl + htmlUrlPath
			longStringUrl = htmlFile
		}
	}
	return longStringUrl
}

func MakeResponseString(toUserName string, fromUserName string, respStr string) string {
	return makeResponseString2(toUserName, fromUserName, "text", respStr)
}
func makeResponseString2(toUserName string, fromUserName string, msgType string, respStr string) string {
	respInfo := domain.WXRespTextMsg{}
	respInfo.FromUserName = domain.CDATA{toUserName}
	respInfo.ToUserName = domain.CDATA{fromUserName}
	respInfo.MsgType = domain.CDATA{msgType}

	//affixString := ""
	//if len(respStr) >= max_length_wechat {
	//	affixString = "[...]"
	//}
	respInfo.Content = domain.CDATA{utils.SubstringByBytesWholeChar(respStr, max_length_wechat)}
	respInfo.CreateTime = time.Now().Unix()

	respXml2String, _ := xml.MarshalIndent(respInfo, "", "")
	return string(respXml2String)
}
