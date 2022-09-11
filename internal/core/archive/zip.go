package archive

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"os"

	"filetransfer-back/internal/config"
)

type ZipArchiver struct {}

func NewZipArchiver () *ZipArchiver {
  za := &ZipArchiver{}
  return za
}
func (za *ZipArchiver) Extension() string {
  return config.ZipExtension()
}
func (za *ZipArchiver) Archive (efiles []EntryFile) ([]byte, error) {
  var err error
  var zf io.Writer 
  var buf bytes.Buffer
  
  sep := string(os.PathSeparator)
  zWriter := zip.NewWriter(&buf)
  for _, efile := range efiles {
    fi := efile.FileInfo()
    if fi != nil {
      header, err := zip.FileInfoHeader(fi)
      if err == nil {
        header.Method = zip.Deflate
        if fi.IsDir() {
          header.Name += sep
        }
        zf, err = zWriter.CreateHeader(header)
        if err == nil {
          _, err = zf.Write(efile.Data())
          if err != nil {
            break
          }
        } else {
          break
        }
      } else {
        break
      }
    } else {
      err = errors.New("FileInfo cannot be nil")
      break
    }
  }
  err = zWriter.Close()
  return buf.Bytes(), err
}

func (za *ZipArchiver) Unarchive (adata []byte) ([]EntryFile, error) {
  var err error
  var efiles []EntryFile
  var tbuf bytes.Buffer

  bReader := bytes.NewReader(adata)
  zReader, err := zip.NewReader(bReader, int64(len(adata)))
  
  if err == nil {
    for _, f := range zReader.File {
      tbuf.Reset()
      fileInArchive, err := f.Open()
      if err == nil {
        _, err := io.Copy(&tbuf, fileInArchive)
        if err == nil {
          fi := f.FileInfo()
          fileInArchive.Close()
          
          efile, err := NewEntryFile(WithFileInfo(fi), WithData(tbuf.Bytes()))
          if err == nil {
            efiles = append(efiles, *efile)
          } else {
            break
          }
        } else {
          break
        }
      } else {
        break
      }
    }

  }
  
  return efiles, err
}

