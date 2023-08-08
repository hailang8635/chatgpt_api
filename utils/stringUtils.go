package utils

import (
    "encoding/base64"
    "fmt"
    "strings"
)

func SubstringByBytes(str string, length int) string {

    // runeData := []rune(str)
    // length := 100
    if len(str) < length {
        //length = len(str)
        return str
    } else {
        return str[0:length] + "..."
    }

}


func Substring(str string, length int) string {
    runeData := []rune(str)
    // length := 100
    if len(runeData) < length {
        //length = len(runeData)
        return str
    } else {
        return string(runeData[0:length]) + "..."
    }

}

func Base64Encode(msgStr string) string {
    return base64.StdEncoding.EncodeToString([] byte(msgStr))

}

func Base64Decode(msg string) string{

    msg = strings.ReplaceAll(msg, "\r", "")
    msg = strings.ReplaceAll(msg, "\n", "")
    msg = strings.Trim(msg, " ")


    //base64.StdEncoding
    decode, err := base64.StdEncoding.DecodeString(msg)

    if err != nil {
        return ""
    } else {
        return strings.ReplaceAll(string(decode), "\n", "")
    }

}

func TestBase64() {
    encode := Base64Encode("IiIi5ZibLuaJuQo=")
    fmt.Println("encode: ", encode)
    fmt.Println("decode: ", Base64Decode("6Ieq54Sa"))

    //filePath := "../"
}