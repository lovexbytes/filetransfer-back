package model

import (
  "errors"
  "database/sql"

  "filetransfer-back/internal/config"
  "filetransfer-back/internal/pkg/db"

)

type _ssInfo struct {
  uuid string
  slug string
  ttl int32
}

type SsInfoOption func (*_ssInfo)

func NewSsInfo (opts ...SsInfoOption) (*_ssInfo, error) {
  var (
    err error
    ssInfo *_ssInfo
  )
  
  for _, opt := range opts {
    opt(ssInfo)
  }

  err = ssInfo.validate()

  return ssInfo, err
}

func (ssInfo *_ssInfo) validate () error {

  if len(ssInfo.uuid) == 0 {
    return errors.New("SsInfo uuid must not be empty")
  }
  
  if len(ssInfo.slug) == 0 {
    return errors.New("SsInfo slug must not be empty")
  }

  if ssInfo.ttl < 0 {
    return errors.New("SsInfo ttl must be >= 0")
  }

  return nil
}

func (ssInfo *_ssInfo) Save () error {
  var (
    err error
    conn *sql.DB
  )
  
  conn, err = db.GetDB()

  if err == nil {
    defer conn.Close()

    var stmt *sql.Stmt
    stmt, err = conn.Prepare("INSERT INTO ? (uuid, slug, ttl) VALUES (?, ?, ?)")
    
    if err == nil {
      defer stmt.Close()
      _, err = stmt.Exec(config.DbTableName(), ssInfo.UUID(), ssInfo.Slug(), ssInfo.TTL())
    }

  }

  return err
}

//options
func WithUUID (uuid string) SsInfoOption {
  return func (ssInfo *_ssInfo) {
    ssInfo.uuid = uuid
  }
}

func WithSlug (slug string) SsInfoOption {
  return func (ssInfo *_ssInfo) {
    ssInfo.slug = slug
  }
}

func WithTTL (ttl int32) SsInfoOption {
  return func (ssInfo *_ssInfo) {
    ssInfo.ttl = ttl
  }
}

//getters
func (ssInfo *_ssInfo) UUID() string {
  return ssInfo.uuid
}

func (ssInfo *_ssInfo) Slug() string {
  return ssInfo.slug
}

func (ssInfo *_ssInfo) TTL() int32 {
  return ssInfo.ttl
}
