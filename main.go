// +build !unit

package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// @NOTE: the real connection is not required for tests
	db, err := sql.Open("mysql", "root@/blog")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			panic(err)
		}
	}()

	if err = recordStats(db, 1 /*some user id*/, 5 /*some product id*/); err != nil {
		panic(err)
	}
}
