package db

import (
  "testing"
)

func TestGetDB (t *testing.T) {
  db, err := GetDB()

  if err != nil {
    t.Fatal(err)
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    t.Fatal(err)
  }

}

