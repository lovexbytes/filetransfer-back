package config

import (
)

const (
  //root package path
  _dataFolderPath = "/tmp"

  _ssExtension = "ss"
)

func DataFolderPath () string {
  /**
    @TODO
    get it from database or something
  */
  return _dataFolderPath
}

func SecureStorageExtension () string {
  return _ssExtension
}

