package dbutils

import (
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)



type UserInfo struct {
    ID     uint
    Name   string
    Gender string
    Hobby  string
}


func Test() {
    //dsn := "root:@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
    dsn := "chatgpt:your_password@tcp(127.0.0.1:3306)/chatgpt?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic(err)
    }
    // 创建表自动迁移, 把结构体和数据表进行对应
    db.AutoMigrate(&UserInfo{})

    // 创建数据信息
    createUserinfo := UserInfo{1, "randyRandy", "man", "code"}
    db.Create(&createUserinfo)

    // 查询
    var u UserInfo
    db.First(&u)
    fmt.Printf("userInfo:%#v", u)

    // 更新
    db.Model(&u).Update("hobby", "Python Go code")
    fmt.Printf("userInfo:%#v", u)

    // 删除
    db.Delete(&u)

}
