package data

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/stewartwebb/filestore/src/common"
)

type DB struct {
	*sql.DB
}

func NewDB() (*DB, error) {
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		common.AppConfig.DatabaseUser, common.AppConfig.DatabasePassword,
		common.AppConfig.DatabaseHost, common.AppConfig.DatabasePort,
		common.AppConfig.DatabaseName)
	db, err := sql.Open("mysql", dbinfo)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// AttachDatabase adds the database connection to the http context
func AttachDatabase(next httprouter.Handle, db *DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "conn", db)
		next(w, r.WithContext(ctx), ps)
	}
}

// AttachDatabaseOld adds the database connection to the http context
func AttachDatabaseOld(db *DB) (mw func(http.Handler) http.Handler) {
	mw = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "conn", db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	return
}
