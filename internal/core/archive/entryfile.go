package archive

import (
  "io/fs"
  "errors"
)

type EntryFile struct {
  fileInfo fs.FileInfo
  data []byte
}

type EntryFileOption func (*EntryFile)

func NewEntryFile(opts ...EntryFileOption) (*EntryFile, error) {
  ef := &EntryFile{
    fileInfo: nil,
    data : nil,
  }

  for _, opt := range opts {
    opt(ef)
  }
  
  err := ef.validate()
  if err != nil {
    ef = nil
  }
  return ef, err
}

func (ef *EntryFile) validate () error {
  var err error 
  if ef.fileInfo != nil {
    if ef.data != nil {
      if len(ef.data) > 0 {} else {
        err = errors.New("data cannot be empty")
      }
    } else {
      err = errors.New("data cannot be nil")
    }
  } else {
    err = errors.New("fileInfo cannot be nil")
  }
  return err
}

//getters

func (ef *EntryFile) FileInfo() fs.FileInfo{
  return ef.fileInfo
}

func (ef *EntryFile) Data() []byte {
  return ef.data
}

//options

func WithFileInfo(fi fs.FileInfo) EntryFileOption {
  return func (ef *EntryFile) {
    ef.fileInfo = fi
  }
}

func WithData (data []byte) EntryFileOption {
  return func (ef *EntryFile) {
    ef.data = data
  }
}




