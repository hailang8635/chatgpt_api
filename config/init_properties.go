package config

import (
    "chatgpt_api/utils"
    "fmt"
    "github.com/gobuffalo/packr"
    "log"
    "strings"
)

var bad_words []string

func init() {

    text := getPackrText(".", "ban_words.txt")
    text += getPackrText("./config/", "ban_words.txt")

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
    log.Println(box.List(), box.Path)
    text, err := box.FindString(filename)

    if err != nil {
        log.Println("init_properties 1", err)
    } else {
        log.Println("init_properties 2", utils.Substring(text, 10))
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
            //log.Println("wordStr : ", wordStr, " badWord: ", badWord)
            wordStr = strings.ReplaceAll(wordStr, badWord, "**")
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
