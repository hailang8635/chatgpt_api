package batch_loop_ask_task

import (
	"chatgpt_api/utils"
	"encoding/csv"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"os"
	"strings"
)

func ReadCsv(fileName string) string {
	// 打开CSV文件（请将"input.csv"替换为你的实际文件路径）
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}
	defer file.Close() // 确保文件在函数结束时关闭

	// 转换为UTF-8阅读器
	fileUTF8 := transform.NewReader(file, simplifiedchinese.GBK.NewDecoder())

	// 创建CSV Reader
	reader := csv.NewReader(fileUTF8)

	// 可选配置（根据实际CSV格式调整）
	reader.Comma = ','       // 设置分隔符，默认为逗号
	reader.Comment = '#'     // 设置注释标识符
	reader.LazyQuotes = true // 允许非标准引号格式

	// 记录处理的行数
	lineCount := 0
	var titleLine []string

	for {
		// 读取一行记录
		record, err := reader.Read()
		if err == io.EOF {
			break // 文件读取完成
		}
		if err != nil {
			// 处理读取错误但不终止程序（可根据需求调整）
			log.Printf("第 %d 行解析错误: %v", lineCount+1, err)
			continue
		}

		lineCount++

		// 处理记录（这里简单打印，可根据需求修改）
		//processRecord(lineCount, record)
		//return record, lineCount

		if lineCount == 1 {
			titleLine = record
		} else {
			respStr := Ask_single(titleLine, record)
			utils.AppendFile("./ai_output_"+fileName+".txt", respStr+"\n")
			//, reqStr+"\t\t"+respStr+"\n")
		}

	}

	fmt.Printf("读取 %d 行数据\n", lineCount)
	return "处理完毕行数：" + string(lineCount)
}

// 处理单行记录的函数
func MakeRecordToString(titleLine []string, record []string) string {
	// 校验列数是否匹配
	if len(titleLine) != len(record) {
		//panic("列数不匹配")
		fmt.Println("列数不匹配", record[0], record[1])
		return ""
	}

	// 构造结果字符串
	var builder strings.Builder
	for i := 0; i < len(titleLine); i++ {
		builder.WriteString(fmt.Sprintf("%s:%s", titleLine[i], record[i]))
		if i < len(titleLine)-1 {
			builder.WriteString(", ")
		}
	}

	return builder.String()
}
func processRecord(lineNumber int, record []string) {
	fmt.Printf("第 %d 行: ", lineNumber)
	for i, field := range record {
		fmt.Printf("字段%d[%s] ", i+1, field)
	}
	fmt.Println()
}
