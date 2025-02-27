package api_from_ai

import (
	"chatgpt_api/config"
	"chatgpt_api/domain"
	"chatgpt_api/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
//API_KEY = ""
//API_URL = "https://open.bigmodel.cn/api/paas/v4/chat/completions/"
)

/**
 * 返回答案，字符串的格式
 *
 */
func DeepSeekApi(content string) (string, error) {
	return DeepSeekApiWithHistory(content, nil)
}

func DeepSeekApiWithHistory(content string, keywordsArr []domain.KeywordAndAnswerItem) (string, error) {
	log.Printf("  --> ask deepseek.com ...【%s】 (%s)\n", content, config.DeepseekUrl)

	messagesInfo := []Message{}

	for _, keywords := range keywordsArr {
		messagesInfo = append(messagesInfo, Message{Role: "user", Content: keywords.Keyword})
		// messagesInfo = append(messagesInfo, Message{Role: "system", Content: keywords.Answer})
	}
	log.Printf("附带历史消息 %d 条", len(keywordsArr))

	messagesInfo = append(messagesInfo, Message{Role: "user", Content: content})

	// 构建请求体
	// 你好，请用一句话介绍你自己
	requestBody := DeepseekRequest{
		// deepseek-chat 模型已全面升级为 DeepSeek-V3，接口不变。 通过指定 model='deepseek-chat' 即可调用 DeepSeek-V3。
		Model:    config.DeepseekModel,
		Messages: messagesInfo,
		//Messages: []Message{{Role: "user", Content: content}},
	}

	jsonBody, _ := json.Marshal(requestBody)

	//fmt.Println("jsonBody request", string(jsonBody))

	// 创建 HTTP 请求
	//req, _ := http.NewRequest("POST", config.DeepseekUrl, bytes.NewBuffer(jsonBody))
	req, _ := http.NewRequest("POST", config.DeepseekUrl, strings.NewReader(string(jsonBody)))
	req.Header.Set("Authorization", "Bearer "+config.DeepseekApiKey)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("  --> req json: %s %s", config.DeepseekUrl, string(jsonBody))

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 处理响应
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Println("StatusCode错误 resp.Status:", resp.Status)
	}

	var result DeepseekResponse
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}

	// fmt.Println("DeepseekResponse: ", result)
	if result.Error.Code == "" {
		fmt.Println("  <-- 回复:", utils.Substring(result.Choices[0].Message.Content, 100))

		respContentStr := result.Choices[0].Message.Content
		log.Printf("  <-- %s 调用 deepseek.com 完成 %s", content, utils.Substring(respContentStr, 20))
		return respContentStr, nil
	} else {
		fmt.Println("  <-- 错误:", result.Error.Message)
		//return "", GptApiError(resp.StatusCode, resp.Status, result.Error.Message)
		return "", nil
	}

	//respStr := string(respBody)

	//jsonReader, _ := simplejson.NewFromReader(strings.NewReader(respStr))
	//respContentStr, _ := jsonReader.Get("choices").GetIndex(0).Get("message").Get("content").String()

}

type DeepseekResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type DeepseekRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type DeepseekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

/**
{
  "error": {
    "message": "Authentication Fails (no such user)",
    "type": "authentication_error",
    "param": null,
    "code": "invalid_request_error"
  }
}
{
  "id": "930c60df-bf64-41c9-a88e-3ec75f81e00e",
  "choices": [
    {
      "finish_reason": "stop",
      "index": 0,
      "message": {
        "content": "Hello! How can I help you today?",
        "role": "assistant"
      }
    }
  ],
  "created": 1705651092,
  "model": "deepseek-chat",
  "object": "chat.completion",
  "usage": {
    "completion_tokens": 10,
    "prompt_tokens": 16,
    "total_tokens": 26
  }
}
*/
