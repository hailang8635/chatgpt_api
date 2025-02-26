package wechat_server

import (
	"chatgpt_api/api_from_ai"
	"chatgpt_api/config"
	"chatgpt_api/domain"
	"fmt"
	"log"
)

/**
 * api string,
 */
func GetAPIResult(content string, item []domain.KeywordAndAnswerItem) (string, error) {
	// Mock TODO
	api := config.DefaultAPI

	if config.SwitchForMockOfAiApi {
		log.Println("使用mock接口，不真实调用外部api，SwitchForMockOfAiApi: ", config.SwitchForMockOfAiApi)
		return config.ApiResponseString, nil
	}

	if api == "gpt" {
		return api_from_ai.GptApi2(content, item)

	} else if api == "glm" {
		return api_from_ai.GLMApiWithHistory(content, item)

	} else {
		return api_from_ai.GLMApiWithHistory(content, item)
	}

}

// TODO
type ActionFunc func() (interface{}, error)

// 创建一个映射，将字符串映射到对应的函数
var actionMap = map[string]ActionFunc{
	"action1": action1,
	"action2": action2,
	// ... 可以继续添加其他动作和对应的函数
}

// 示例动作函数1
func action1() (interface{}, error) {
	return "Executing action 1", nil
}

// 示例动作函数2
func action2() (interface{}, error) {
	return 42, nil // 假设返回了一个整数结果
}

// 执行动作的函数，根据传入的动作名称调用相应的函数
func ExecuteAction(actionName string) (interface{}, error) {
	// 从映射中查找对应的函数
	action, exists := actionMap[actionName]
	if !exists {
		//return nil, errors.New("unknown action")
		return nil, nil
	}
	// 调用找到的函数
	return action()
}

func mainabc() {
	// 示例调用
	actionName := "action1"
	result, err := ExecuteAction(actionName)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
}
