package database

import (
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"os"
)

var mutexMainDatabase sync.Mutex
var wankeliaoDatabase *sql.DB
var ErrNoRows = sql.ErrNoRows

func Linkdb() {
	var err error
	wankeliaoDatabase, err = sql.Open("mysql", os.Getenv("dbAccount")+":"+os.Getenv("dbPassword")+"@tcp("+os.Getenv("dbHost")+")/wankeliao")
	for err != nil {
		log.Printf("Linkdb err %+v\n", err)
		now := time.Now().UnixNano()
		for time.Now().UnixNano() <= now+1e9 {
		}
		wankeliaoDatabase, err = sql.Open("mysql", os.Getenv("dbAccount")+":"+os.Getenv("dbPassword")+"@tcp("+os.Getenv("dbHost")+")/wankeliao")
	}
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	mutexMainDatabase.Lock()
	defer mutexMainDatabase.Unlock()
	return wankeliaoDatabase.Exec(query, args...)
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	mutexMainDatabase.Lock()
	defer mutexMainDatabase.Unlock()
	return wankeliaoDatabase.Query(query, args...)
}

func QueryRow(query string, args ...interface{}) *sql.Row {
	mutexMainDatabase.Lock()
	defer mutexMainDatabase.Unlock()
	return wankeliaoDatabase.QueryRow(query, args...)
}
