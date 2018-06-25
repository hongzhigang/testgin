package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	//数据库连接
	var err error
	db, err = sql.Open("mysql", "root:1234@tcp(192.168.6.200:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	//defer db.Close()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
}

type User struct {
	Id       string `form:"id" json:"id"`
	Page     string `form:"page" json:"page"`
	Username string `form:"username" json:"username"`
	Password int    `form:"password" json:"password"`
}

//点赞
type Like struct {
	Uid    uint `json:"uid"`
	Tid    uint `json:"tid"`
	Cancel bool `json:cancel`
}

//点赞了（可以把DBWriteLike作为Like的方法）
func DBWriteLike(l *Like) error {
	//写法1
	// stmt, err := db.Prepare("INSERT INTO t_like (uid, tid, cancel) VALUES(?,?,?)")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// result, err := stmt.Exec(l.Uid, l.Tid, l.Cancel)
	_, err := db.Exec("INSERT INTO t_like (uid, tid, cancel) VALUES(?,?,?)", l.Uid, l.Tid, l.Cancel)

	return err
}

//获取点赞信息（可以把DBReadLikeByUid作为Like的方法）
func DBReadLikeByUid(uid uint) (titles []uint) {
	//select tid from t_like where uid = 123;
	var autocommit string
	if err := db.QueryRow("SELECT @@autocommit").Scan(&autocommit); err != nil {
		log.Fatal(err)
	}
	fmt.Println(autocommit)

	rows, err := db.Query("SELECT tid FROM t_like WHERE uid = ?", uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	titles = make([]uint, 0)
	for rows.Next() {
		var strTid string
		if err := rows.Scan(&strTid); err != nil {
			log.Fatal(err)
		}
		tid, _ := strconv.ParseUint(strTid, 10, 64)
		titles = append(titles, uint(tid))
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return
}