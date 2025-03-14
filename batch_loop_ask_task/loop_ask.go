package batch_loop_ask_task

import (
	"chatgpt_api/config"
	"chatgpt_api/wechat_server"
	"fmt"
)

func Ask_batch(fileName string) string {

	resultMsg := ReadCsv(fileName)
	fmt.Println(resultMsg)

	return resultMsg
}

//
func Ask_single(titleLine []string, record []string) string {

	// 前缀、后缀配置化
	// 请帮我判断一下, 中国东方航空股份有限公司, 基于国民经济行业分类（GB/T 4754-2017)标准下钻到四级分类, 仅返回json格式分类
	// 请帮我判断一下, 中国东方航空股份有限公司, 基于国民经济行业分类（GB/T 4754-2017)标准下钻到四级分类, 不加说明仅返回tab分割的一、二、三、四级分类名称
	//keywordNew := "请帮我判断一下, " + keyword + ", 基于国民经济行业分类（GB/T 4754-2017)标准下钻到四级分类, 不加说明仅返回|分割的一、二、三、四级分类名称"
	reqStr := MakeRecordToString(titleLine, record)
	keywordNew := config.Batch_ask_prefix + reqStr + config.Batch_ask_suffix
	// 根据关键词查询GPT接口
	respStr, err := wechat_server.GetAPIResult(keywordNew, nil)
	if err != nil {
		fmt.Println("系统忙，请稍后再试.", err)
		return ""
	} else {
		fmt.Println("Ask_single完成 ", reqStr, respStr)
		return reqStr + "|" + respStr
	}
}
