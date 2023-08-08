package domain

import (
    "encoding/xml"
    "time"
)

type CDATA struct {
    Text string `xml:",cdata"`
}

/*
 *
    1:是，2:否，默认0

    Is_done      GPT是否已查得
    Is_finished  是否已返回给用户

*/

type Keywords struct {
    //gorm.Model
    Id          int64     `gorm:"id"`
    Catalog     string    `gorm:"catalog"`
    Fromuser    string    `gorm:"fromuser"`
    Keyword     string    `gorm:"keyword"`
    Answer      string    `gorm:"answer"`
    Labels      string    `gorm:"labels"`
    Is_finished int       `gorm:"is_finished"`
    Is_done     int       `gorm:"is_done"`
    Create_time time.Time `gorm:"create_time"`
    Finish_time time.Time `gorm:"finish_time"`

    Create_time_start time.Time `gorm:"-"`
    Create_time_end   time.Time `gorm:"-"`
    OrderByIdDesc     bool      `gorm:"-"`
}

// WXReqTextMsg 微信文本消息结构体
type WXReqTextMsg struct {
    ToUserName   string
    FromUserName string
    CreateTime   int64
    MsgType      string
    Content      string
    MediaId      string // 语音消息媒体id，可以调用获取临时素材接口拉取该媒体
    Format       string // 语音格式：amr
    Recognition  string
    MsgId        int64
    MsgDataId    int64
    Idx          int64
}

// WXRespTextMsg 微信回复文本消息结构体
type WXRespTextMsg struct {
    ToUserName   CDATA
    FromUserName CDATA
    CreateTime   int64
    MsgType      CDATA
    Content      CDATA

    // 若不标记XMLName, 则解析后的xml名为该结构体的名称
    XMLName xml.Name `xml:"xml"`
}

type RespMsg struct {
    Keyword    string
    Username   string
    Answer     string
    IsDone     int
    IsReturned int
}
