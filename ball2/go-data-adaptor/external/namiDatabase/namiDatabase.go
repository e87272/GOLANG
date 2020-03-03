package namiDatabase

import (
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"os"
)

var mutexMainDatabase = new(sync.Mutex)
var namiDatabase *sql.DB
var ErrNoRows = sql.ErrNoRows

func Linkdb() {
	var err error
	namiDatabase, err = sql.Open("mysql", os.Getenv("dbAccount")+":"+os.Getenv("dbPassword")+"@tcp("+os.Getenv("dbHost")+")/nami")
	for err != nil {
		// log.Printf("Linkdb err %+v\n", err)
		now := time.Now().UnixNano()
		for time.Now().UnixNano() <= now+1e9 {
		}
		namiDatabase, err = sql.Open("mysql", os.Getenv("dbAccount")+":"+os.Getenv("dbPassword")+"@tcp("+os.Getenv("dbHost")+")/nami")
	}
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	mutexMainDatabase.Lock()
	defer mutexMainDatabase.Unlock()
	return namiDatabase.Exec(query, args...)
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	mutexMainDatabase.Lock()
	defer mutexMainDatabase.Unlock()
	return namiDatabase.Query(query, args...)
}

func QueryRow(query string, args ...interface{}) *sql.Row {
	mutexMainDatabase.Lock()
	defer mutexMainDatabase.Unlock()
	return namiDatabase.QueryRow(query, args...)
}

func Begin() (*sql.Tx, error) {
	mutexMainDatabase.Lock()
	defer mutexMainDatabase.Unlock()
	return namiDatabase.Begin()
}
func Prepare(tx *sql.Tx, query string) (*sql.Stmt, error) {
	mutexMainDatabase.Lock()
	defer mutexMainDatabase.Unlock()
	return namiDatabase.Prepare(query)
}

func Commit(tx *sql.Tx) error {
	mutexMainDatabase.Lock()
	defer mutexMainDatabase.Unlock()
	return tx.Commit()
}
