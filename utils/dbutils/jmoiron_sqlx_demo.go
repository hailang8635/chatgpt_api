package dbutils

import (
    "fmt"
    "chatgpt_api/domain"
    "reflect"
    "strings"

    //"database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
    "log"
    "time"
)

/*
type Keywords struct {
    Id           int64      `db:"id"`
    Catalog      string     `db:"catalog"`
    Fromuser     string     `db:"fromuser"`
    Keyword      string     `db:"keyword"`
    Answer       string     `db:"answer"`
    Labels       string     `db:"labels"`
    Is_finished  int        `db:"is_finished"`
    Is_done      int        `db:"is_done"`
    Create_time  time.Time  `db:"create_time"`
    Finish_time  time.Time  `db:"finish_time"`
}
*/

var db *sqlx.DB

func initdb_sqlx() error {
    db, err := sqlx.Open("mysql", "chatgpt:your_password@(127.0.0.1:3306)/chatgpt?charset=utf8mb4&parsetime=true")
    if err != nil {
        log.Println("open mysql failed,", err)
    }

    err = db.Ping()
    if err != nil {
        log.Println("sql Ping exception", err)
    } else {
        defer db.Close()
    }

    db.SetMaxIdleConns(30)
    db.SetMaxOpenConns(50)

    return err
}

func Insert_sqlx(fromuser string, keyword string, answer string, is_finished int, is_done int, create_time string, finish_time string) {

    err := initdb_sqlx()

    if err != nil {
        log.Println("sql init exception")
        return
    } else {
        defer db.Close()
    }

    insertDML := `INSERT INTO chatgpt.keywords000 (fromuser, keyword, answer, is_finished, is_done, create_time, finish_time) VALUES (?, ?, ?, ?, ?, ?, ?)`

    result, err := db.Exec(insertDML, fromuser, keyword, answer, is_finished, is_done, time.Now(), nil)
    affected, _ := result.RowsAffected()
    id, _ := result.LastInsertId()

    log.Println("affected/id/err: ", affected, id, err)
}

func Select_sqlx(fromuser string, keyword string, answer string, is_finished int, is_done int) []domain.Keywords {
    var keywordsArr []domain.Keywords

    // db, err := sql.Open("mysql", "chatgpt:your_password@(127.0.0.1:3306)/chatgpt?parseTime=true")

    selectDML := `select * from chatgpt.keywords000 where `
    // fromuser, keyword, answer, is_finished, is_done, create_time, finish_time
    // fromuser, keyword, answer, is_finished, is_done, time.Now(), nil)

    queryParams := QueryParams{
        "fromuser":    fromuser,
        "keyword":     keyword,
        "answer":      answer,
        "is_finished": is_finished,
        "is_done":     is_done,
    }

    //selectDML
    whereCondition, values := MakeWhereCondition(queryParams)
    selectDML = selectDML + whereCondition

    log.Println(selectDML, values)

    err := db.Select(&keywordsArr, selectDML)

    //err := db.Query(&keywordsArr, selectDML, values)

    if err != nil {
        log.Println(err)
        return nil
    } else {
        return keywordsArr
    }

}

type QueryParams map[string]interface{}

func MakeWhereCondition(params QueryParams) (string, []interface{}) {
    var values []interface{}
    var where []string
    for k, v := range params {

        if reflect.TypeOf(v).String() == "string" {
            valueString := reflect.ValueOf(v).String()

            if valueString == "" {
                continue
            }
            values = append(values, valueString)
        } else if reflect.TypeOf(v).String() == "int" {
            value := reflect.ValueOf(v).Int()

            if value == 0 {
                continue
            }
            //values = append(values, valueString)
        }

        //MySQL Way:
        where = append(where, fmt.Sprintf(" %s = ? ", k))

        // where = append(where, fmt.Sprintf(`"%s" = %s`,k, "$" + strconv.Itoa(len(values))))
    }

    // string := ("SELECT name FROM users WHERE " + strings.Join(where, " AND "))
    //for testing purposes i didn't ran actual query, just print it in the console and returned JSON back

    whereCondition := strings.Join(where, " AND ")

    // fmt.Println(whereCondition, values)
    return whereCondition, values

}
