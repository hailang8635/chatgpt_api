package wechat_server

import (
	"fmt"
	"net/http"
)

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
		respStr, err := GetAPIResult(s, nil)
		if err != nil {
			fmt.Fprintf(w, "%s %s", s, "系统忙，请稍后再试.")
			return
		}
		fmt.Fprintf(w, "【问】：%s \n\n【答】: \n %s", s, respStr)
	}
}
