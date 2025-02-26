package wechat_server

import (
	"bytes"
	"chatgpt_api/config"
	"chatgpt_api/domain"
	"chatgpt_api/utils"
	"encoding/xml"
	"fmt"
	"github.com/yuin/goldmark"
	"log"
	"os"
	"strings"
	"time"
)

var timeLayoutStrYYYYMMDDHHmmss = "20060102150405"

var keywords = map[string]domain.RespMsg{}

//var retry_gap int64 = 5
//var length_wechat = 500

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

	offset_5m, _ := time.ParseDuration("-5m")
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
	if len(respStr) > length_wechat || strings.Contains(respStr, "```") {
		// TODO
		var buf bytes.Buffer
		err := goldmark.Convert([]byte(respStr), &buf)
		if err != nil {
			log.Println("markdown --> html, exception", err)
		} else {
			log.Println("markdown --> html, ", utils.Substring(buf.String(), 20))

			fileNameRight := utils.Substring(strings.ReplaceAll(keywordParamsOrigin, " ", ""), 12)
			// TODO 同步至okzhang.com
			htmlFile := startTime.Format(timeLayoutStrYYYYMMDDHHmmss) + "_" + fileNameRight + ".html"
			//htmlUrlPath := startTime.Format(timeLayoutStrYYYYMMDDHHmmss) + "_" + url.QueryEscape(fileNameRight) + ".html"
			file, err := os.Create(config.HtmlDir + htmlFile)
			if err != nil {
				fmt.Println("create html file error", err)
			}
			file.WriteString("<html><head>  <title>ChatGPT助手-安德鲁家的550W</title>  <basefont face=\"微软雅黑\" size=\"2\" />  <meta http-equiv=\"Content-Type\" content=\"text/html;charset=utf-8\" /></head>")
			defer file.Close()

			_, err = buf.WriteTo(file)
			if err != nil {
				fmt.Println("write html file error", err)
			}

			// https://chatapi.okzhang.com/html/cah/test.html
			// "[答案详情见链接] \n" + utils.HtmlUrl + htmlUrlPath
			longStringUrl = htmlFile
		}
	}
	return longStringUrl
}

func MakeResponseString(toUserName string, fromUserName string, respStr string) string {
	return MakeResponseString2(toUserName, fromUserName, "text", respStr)
}
func MakeResponseString2(toUserName string, fromUserName string, msgType string, respStr string) string {
	respInfo := domain.WXRespTextMsg{}
	respInfo.FromUserName = domain.CDATA{toUserName}
	respInfo.ToUserName = domain.CDATA{fromUserName}
	respInfo.MsgType = domain.CDATA{msgType}
	respInfo.Content = domain.CDATA{utils.SubstringByBytesWholeChar(respStr, length_wechat)}
	respInfo.CreateTime = time.Now().Unix()

	respXml2String, _ := xml.MarshalIndent(respInfo, "", "")
	return string(respXml2String)
}
