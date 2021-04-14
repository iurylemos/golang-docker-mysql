package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Connectar() (*sql.DB, error) {
	urlConn := "root:password@/devbook?charset=utf8&parseTime=True&loc=Local"

	db, erro := sql.Open("mysql", urlConn)

	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		return nil, erro
	}

	return db, nil
}
