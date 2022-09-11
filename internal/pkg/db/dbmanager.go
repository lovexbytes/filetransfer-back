package db

import (
  "database/sql"
  "errors"

  "filetransfer-back/internal/config"

  _ "github.com/mattn/go-sqlite3"
)

var _conn *sql.DB

func GetDB() (*sql.DB, error) {
  var err error

  if _conn == nil {
    _conn, err = connect()
  } else {
    err = _conn.Ping()
  }

  return _conn, err
}

func connect () (*sql.DB, error) {
  var (
    db *sql.DB
    err error
    dbpath = config.DbPath()
  )
  if len(dbpath) > 0 {
    db, err = sql.Open("sqlite3", dbpath)
    if err != nil {
      err = db.Ping()
    }
  } else {
    err = errors.New("Empty db filepath")
  }
  
  return db, err
}

