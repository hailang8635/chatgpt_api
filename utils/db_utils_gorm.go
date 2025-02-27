package utils

import (
	"chatgpt_api/domain"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	//"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

var GlobalConn *gorm.DB

// 启动时每个包自动执行init()方法
func init() {

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
	host := config.GetString("database.host")
	port := config.GetString("database.port")
	user := config.GetString("database.user")
	password := config.GetString("database.password")
	dbname := config.GetString("database.dbname")

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	//log.Printf("dsn", string(dsn))

	//dsn := "chatgpt:your_password@tcp(ip:port)/chatgpt?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	GlobalConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	sqlDB, err := GlobalConn.DB()
	if err != nil {
		log.Println("connection err", err)
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxIdleTime(600 * time.Second)

	err = sqlDB.Ping()
	if err != nil {
		sqlDB.Close()
		log.Fatal("ping error")
	}

	// defer db.Close()
	// return conn, nil
}

// open api
func GetDB() *gorm.DB {
	_, err := GlobalConn.DB()

	if err != nil {
		log.Println("connect db server failed.")
	}

	//defer sqlDB.Close()

	return GlobalConn
}

// func Insert(fromuser string, keyword string, answer string, is_finished int, is_done int) int64 {
func Insert(keywords domain.KeywordAndAnswerItem) int64 {
	startTime := time.Now().UnixMilli()

	keywords.Create_time = time.Now()
	keywords.Finish_time = time.Now()
	keywords.Update_time = time.Now()

	// .Debug()
	db := GetDB()
	result := db.Table("t_keywords").Create(&keywords)

	// fmt.Println("insert 行数", result.RowsAffected)
	//return result.RowsAffected
	log.Println(">Insert 耗时:", time.Now().UnixMilli()-startTime, "ms, 数量:", result.RowsAffected, ", 关键字", keywords.Keyword)
	return keywords.Id
}

func Update(keywords domain.KeywordAndAnswerItem) int64 {
	//startTime := time.Now().UnixMilli()

	// .Debug()
	db := GetDB()
	result := db.Table("t_keywords").Model(&keywords).Where("id = ?", keywords.Id).UpdateColumns(domain.KeywordAndAnswerItem{
		Answer:      keywords.Answer,
		Url:         keywords.Url,
		Is_finished: keywords.Is_finished,
		Is_done:     keywords.Is_done,
		Finish_time: keywords.Finish_time,
	})

	//log.Println(">Update 耗时:", time.Now().UnixMilli()-startTime, "ms, 数量:", result.RowsAffected, ", 关键字", keywords.Keyword)
	return result.RowsAffected
}

// func Select(fromuser string, keyword string, answer string, is_finished int, is_done int) (int64, KeywordAndAnswerItem) {
func SelectOne(keywords domain.KeywordAndAnswerItem) (int64, domain.KeywordAndAnswerItem) {
	_, arr := SelectList(keywords, 1)
	if arr != nil && len(arr) >= 1 {
		return 1, arr[0]
	} else {
		return 0, domain.KeywordAndAnswerItem{}
	}
}
func SelectList(keywords domain.KeywordAndAnswerItem, nums int) (int64, []domain.KeywordAndAnswerItem) {
	//startTime := time.Now().UnixMilli()

	// 创建数据信息
	var keywords_result []domain.KeywordAndAnswerItem

	//keywords := KeywordAndAnswerItem{}
	params := map[string]interface{}{}

	// {"name": "jinzhu", "age": 20}
	//
	if keywords.Fromuser != "" {
		params["fromuser"] = keywords.Fromuser
	}
	if keywords.Keyword != "" {
		params["keyword"] = keywords.Keyword
	}
	if keywords.Answer != "" {
		params["answer"] = keywords.Answer
	}
	if keywords.Is_done != 0 {
		params["is_done"] = keywords.Is_done
	}
	if keywords.Is_finished != 0 {
		params["is_finished"] = keywords.Is_finished
	}

	// 创建表自动迁移, 把结构体和数据表进行对应
	//db.AutoMigrate(&KeywordAndAnswerItem{})

	// .Debug()
	// .First()
	db := GetDB()
	db = db.Table("t_keywords")
	if !keywords.Create_time_start.IsZero() {
		//fmt.Println("keywords.Create_time_start", keywords.Create_time_start)
		db = db.Where(" create_time >= ? ", keywords.Create_time_start)
	}

	if !keywords.Create_time_end.IsZero() {
		db = db.Where(" create_time < ? ", keywords.Create_time_end)
	}
	if keywords.OrderByIdDesc {
		db = db.Order("id desc ")
	}

	// 记录排序为按ID顺序递增
	result := db.Where(params).Limit(nums).Find(&keywords_result)

	//log.Println(">Select 耗时:", time.Now().UnixMilli()-startTime, "ms, 数量:", result.RowsAffected, ", 关键字: ", keywords.Keyword)
	return result.RowsAffected, keywords_result

}
