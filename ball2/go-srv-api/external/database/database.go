package database

import (
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"os"
)

var mutexMainDatabase = new(sync.Mutex)
var chiouChiouDatabase *sql.DB
var ErrNoRows = sql.ErrNoRows

func Linkdb() {
	var err error
	chiouChiouDatabase, err = sql.Open("mysql", os.Getenv("dbAccount")+":"+os.Getenv("dbPassword")+"@tcp("+os.Getenv("dbHost")+")/chiou_chiou")
	for err != nil {
		// log.Printf("Linkdb err %+v\n", err)
		timer := time.NewTimer(time.Second)
		select {
		case <-timer.C:
			chiouChiouDatabase, err = sql.Open("mysql", os.Getenv("dbAccount")+":"+os.Getenv("dbPassword")+"@tcp("+os.Getenv("dbHost")+")/chiou_chiou")
		}
	}
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	mutexMainDatabase.Lock()
	defer mutexMainDatabase.Unlock()
	return chiouChiouDatabase.Exec(query, args...)
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	mutexMainDatabase.Lock()
	defer mutexMainDatabase.Unlock()
	return chiouChiouDatabase.Query(query, args...)
}

func QueryRow(query string, args ...interface{}) *sql.Row {
	mutexMainDatabase.Lock()
	defer mutexMainDatabase.Unlock()
	return chiouChiouDatabase.QueryRow(query, args...)
}
