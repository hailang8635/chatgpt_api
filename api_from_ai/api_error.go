package api_from_ai

import "fmt"

type ApiError struct {
	Code    int    // 错误码
	Message string // 错误信息
	Details string // 附加信息
}

// 实现 error 接口
func (e *ApiError) Error() string {
	return fmt.Sprintf("错误 %d: %s (%s)", e.Code, e.Message, e.Details)
}

// 构造函数
func GptApiError(code int, message, details string) *ApiError {
	return &ApiError{
		Code:    code,
		Message: message,
		Details: details,
	}
}
