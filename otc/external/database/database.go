package database

import (
	"database/sql"
	"os"
	"server/common"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var mutexMemberDatabase = new(sync.Mutex)
var memberDatabase *sql.DB
var ErrNoRows = sql.ErrNoRows

func Linkdb() {
	var err error
	for {
		memberDatabase, err = sql.Open("mysql", os.Getenv("dbAccount")+":"+os.Getenv("dbPassword")+"@tcp("+os.Getenv("dbHost")+")/"+os.Getenv("dbName")+"?parseTime=true")
		if err == nil {
			err = memberDatabase.Ping()
		}
		if err == nil {
			break
		}
		common.SysErrorLog(map[string]interface{}{
			"name": "Linkdb member err",
		}, err)
		time.Sleep(time.Second)
	}

	maxLifetime, err := strconv.Atoi(os.Getenv("maxLifetime"))
	if err != nil {
		name := "maxLifetime config err"
		common.SysErrorLog(map[string]interface{}{
			"name": name,
		}, nil)
		panic(err)
	}
	maxOpenConns, err := strconv.Atoi(os.Getenv("maxOpenConns"))
	if err != nil {
		name := "maxOpenConns config err"
		common.SysErrorLog(map[string]interface{}{
			"name": name,
		}, nil)
		panic(err)
	}
	maxIdleConns, err := strconv.Atoi(os.Getenv("maxIdleConns"))
	if err != nil {
		name := "maxIdleConns config err"
		common.SysErrorLog(map[string]interface{}{
			"name": name,
		}, nil)
		panic(err)
	}

	// See "Important settings" section.
	memberDatabase.SetConnMaxLifetime(time.Second * time.Duration(maxLifetime))
	memberDatabase.SetMaxOpenConns(maxOpenConns)
	memberDatabase.SetMaxIdleConns(maxIdleConns)

	common.SysLog(map[string]interface{}{
		"name": "Linkdb " + os.Getenv("dbName") + " ok",
	})
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	mutexMemberDatabase.Lock()
	defer mutexMemberDatabase.Unlock()
	return memberDatabase.Exec(query, args...)
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	mutexMemberDatabase.Lock()
	defer mutexMemberDatabase.Unlock()
	return memberDatabase.Query(query, args...)
}

func QueryRow(query string, args ...interface{}) *sql.Row {
	mutexMemberDatabase.Lock()
	defer mutexMemberDatabase.Unlock()
	return memberDatabase.QueryRow(query, args...)
}
