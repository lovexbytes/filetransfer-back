package compression

import (
  "filetransfer-back/internal/core/archive"
)

type Compressor interface{
  Compress (output string, efiles []archive.EntryFile) error 
  Decompress (outFPath string, filename string) error 

  //utils
  Extension () string
}



