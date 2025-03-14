package utils

import (
	"log"
	"os"
)

func AppendFile(filename string, content string) {
	// 要追加的字符串内容
	//content := "This text will be appended to the file\n"

	// 文件路径（可以是相对或绝对路径）
	//filename := "example.log"

	// 以追加模式打开文件（如果文件不存在则创建）
	// os.O_APPEND - 追加模式
	// os.O_WRONLY - 只写模式
	// os.O_CREATE - 如果文件不存在则创建
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}
	defer file.Close() // 确保文件正确关闭

	// 写入内容到文件
	//bytesWritten, err := file.WriteString(content)
	_, err = file.WriteString(content)
	if err != nil {
		log.Fatalf("写入失败: %v", err)
	}

	// 强制将缓冲区数据写入磁盘（可选）
	err = file.Sync()
	if err != nil {
		log.Printf("警告: 同步文件失败: %v", err)
	}

	// fmt.Printf("成功追加 %d 字节到 %s\n", bytesWritten, filename)
}
