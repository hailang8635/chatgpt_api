package config

import (
	"chatgpt_api/utils"
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/spf13/viper"
	"log"
	"strings"
)

var bad_words []string

var DefaultAPI string
var DefaultPort string

var GptUrl string
var ApiKey string
var ModelVersion string

var GLM_Apikey string
var GLM_Model string
var GLM_Url string

var DeepseekApiKey string
var DeepseekModel string
var DeepseekUrl string

var HtmlDir string
var HtmlUrl string

var SwitchForMockOfAiApi bool
var ApiResponseString string
var Batch_ask_prefix string
var Batch_ask_suffix string

// 启动时每个包自动执行init()方法
//func init() {
func InitProperties() {

	initBadWords()
	initProperties()

}

func initProperties() {
	config := viper.New()
	config.AddConfigPath("./config/")
	config.SetConfigName("gdbc")
	config.SetConfigType("yaml")
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("找不到配置文件.. ./config/gdbc.yaml")
		} else {
			log.Fatal("配置文件出错..")
		}
	}

	// defaultAPI
	DefaultAPI = config.GetString("defaultAPI")
	DefaultPort = config.GetString("defaultPort")
	if DefaultPort == "" {
		DefaultPort = ":80"
	}

	// chatgpt
	ApiKey = config.GetString("apikey")
	ModelVersion = config.GetString("modelVersion")
	GptUrl = config.GetString("gptUrl")

	// glm
	GLM_Apikey = config.GetString("glm.apikey")
	GLM_Model = config.GetString("glm.model")
	GLM_Url = config.GetString("glm.url")

	DeepseekApiKey = config.GetString("deepseek.apikey")
	DeepseekModel = config.GetString("deepseek.model")
	DeepseekUrl = config.GetString("deepseek.url")

	HtmlUrl = config.GetString("htmlUrl")
	HtmlDir = config.GetString("htmlDir")

	SwitchForMockOfAiApi = config.GetBool("switchForMockOfAiApi")
	Batch_ask_prefix = config.GetString("batch_ask_prefix")
	Batch_ask_suffix = config.GetString("batch_ask_suffix")
}

/**
 * 初始化敏感词
 */
func initBadWords() {
	text := getPackrText(".", "ban_words.txt")
	//text += getPackrText("./config/", "ban_words.txt")
	//text += getPackrText("/mnt/e/logs/config", "ban_words.txt")
	text += getPackrText("/opt/chatgpt_api/config", "ban_words.txt")

	lineArr := strings.Split(text, "\n")
	for _, line := range lineArr {

		decode := utils.Base64Decode(line)
		//log.Println(line, decode)
		bad_words = append(bad_words, strings.ToLower(strings.ReplaceAll(decode, "1", "")))
		// fmt.Println(line , "---->" , utils.Base64Decode(line))

	}

	bad_words = append(bad_words, "jjjj")

	log.Println("bad_words len:", len(bad_words))

}

func getPackrText(path string, filename string) string {
	box := packr.NewBox(path)
	log.Println("getPackrText ", filename, "in", box.List(), ", path:[", box.Path, "]")
	text, err := box.FindString(filename)

	if err != nil {
		log.Println("getPackrText 1", err)
	} else {
		log.Println("getPackrText 2", utils.Substring(text, 10))
	}

	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = strings.ReplaceAll(text, "\n\n", "\n")
	return text
}

func VerfiyBadWordsOnlyResult(searchWord string) bool {
	result, _ := VerfiyBadWords(searchWord)
	return result
}

func VerfiyBadWords(wordStr string) (bool, string) {
	result := false
	for _, badWord := range bad_words {
		if wordStr != "" && badWord != "" && strings.Contains(strings.ToLower(wordStr), badWord) || strings.Contains(badWord, strings.ToLower(wordStr)) {
			result = true
			log.Println("wordStr: ", wordStr, ", badWord: ", badWord)
			wordStr = strings.ReplaceAll(wordStr, badWord, "口")
		}
	}
	return result, wordStr
}

func Test() {
	fmt.Println("VerfiyBadWords sb --> ", VerfiyBadWordsOnlyResult("sb"))
	fmt.Println("VerfiyBadWords jjjj --> ", VerfiyBadWordsOnlyResult("jjjj"))
	fmt.Println("VerfiyBadWords j --> ", VerfiyBadWordsOnlyResult("j"))
	fmt.Println("VerfiyBadWords tmd --> ", VerfiyBadWordsOnlyResult("tmd"))
	fmt.Println("VerfiyBadWords 1 --> ", VerfiyBadWordsOnlyResult("1"))
	fmt.Println("VerfiyBadWords cba --> ", VerfiyBadWordsOnlyResult("cba"))
	fmt.Println("VerfiyBadWords 白日依山尽 --> ", VerfiyBadWordsOnlyResult("白日依山尽"))

}
