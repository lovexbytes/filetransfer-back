package config

import (
    "path/filepath"
    "runtime"
    "os"
)

var (
  sep = string(os.PathSeparator)
  _, b, _, _ = runtime.Caller(0)
  _basepath   = filepath.Dir(b)
)

const (
  _dbTableName = "ft_data"
  _dbFileName = "system_data.db"
)

func DbTableName () string {
  return _dbTableName
}

func DbPath () string {
    var (
        respath = ""
        updir = ".."
        relpath = _basepath + sep + updir + sep + _dbFileName
        path, err = filepath.Abs(relpath)
    )

    if err == nil {
        respath = path
    }
    return respath

}
