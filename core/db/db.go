package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

type DB interface {
	Select(string, ...interface{}) (*sql.Rows, error)
	Close()
}

type dbQueryExecutor struct {
	db *sql.DB
}

func (exec *dbQueryExecutor) Select(query string, params ...interface{}) (*sql.Rows, error) {
	p, err := exec.db.Prepare(query)
	if err != nil {
		return &sql.Rows{}, err
	}
	return p.Query(params...)
}

func (exec *dbQueryExecutor) Close() {
	if err := exec.db.Close(); err != nil {
		log.Fatal(err)
	}
}

const (
	HOST     = "devel-postgre.tkpd"
	PORT     = 5432
	USER     = "tkpdtraining"
	PASSWORD = "trainingyangbeneryah"
	DBNAME   = "tokopedia-dev-db"
)

var lock = &sync.Mutex{}
var instance DB

func New() DB {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, USER, PASSWORD, DBNAME))
		fmt.Println(err)
		instance = &dbQueryExecutor{
			db: db,
		}
	}
	return instance
}
