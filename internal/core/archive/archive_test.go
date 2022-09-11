package archive

import (
  "testing"
  "os"
  "io"
  "fmt"
  "errors"
)

func TestZipUnzip(t *testing.T) {

  //init test efiles
  efiles, err := initTestEfiles(t)
  if err != nil {
    t.Fatal(err)
  }

  zipArchiver := NewZipArchiver()
  archiveData, err := zipArchiver.Archive(efiles) 
  if err != nil {
    t.Fatal(err)
  }
  //create file in system
  // archiveFile, err := os.Create("result.zip")
  // if err != nil {
  //   t.Fatal(err)
  // }
  //
  // _, err = archiveFile.Write(archiveData)
  // if err != nil {
  //   t.Fatal(err)
  // }
  // archiveFile.Close()

  efiles, err = zipArchiver.Unarchive(archiveData)
  if err != nil {
    t.Fatal(err)
  }

  if len(efiles) == 0 {
    t.Fatal(errors.New("efiles collection is empty"))
  }

  for _, efile := range efiles {
    if len(efile.Data()) == 0 {
      t.Fatal(errors.New("efile " + efile.FileInfo().Name() + " data is empty"))
    }
    fmt.Println("efile:")
    fmt.Println("  name:       ", efile.FileInfo().Name())
    fmt.Println("  data length:", len(efile.Data()))
    fmt.Println("  data:       ", string(efile.Data()))
  }

}


//util
func initTestEfiles (t *testing.T) ([]EntryFile, error) {
  efiles := []EntryFile{}

  file1, err := os.Create("test1.txt")
  fname1 := file1.Name()

  if err != nil {
    t.Fatal(err)
  }

  file2, err := os.Create("test2.txt")
  fname2 := file2.Name()
  if err != nil {
    t.Fatal(err)
  }

  _, err = file1.WriteString("test1 string")
  if err != nil {
    t.Fatal(err)
  }

  _, err = file2.WriteString("test2 string")
  if err != nil {
    t.Fatal(err)
  }

  file1.Close()
  file2.Close()
  
  file1, err = os.Open(fname1)
  if err != nil {
    t.Fatal(err)
  }
  file2, err = os.Open(fname2)
  if err != nil {
    t.Fatal(err)
  }

  fi1, err := file1.Stat()
  if err != nil {
    t.Fatal(err)
  }

  fi1Data, err := io.ReadAll(file1)
  if err != nil {
    t.Fatal(err)
  }
  file1.Close()

  efile1, err := NewEntryFile(WithFileInfo(fi1), WithData(fi1Data))
  if err != nil {
    t.Fatal(err)
  }

  fi2, err := file2.Stat()
  if err != nil {
    t.Fatal(err)
  }

  fi2Data, err := io.ReadAll(file2)
  if err != nil {
    t.Fatal(err)
  }

  file2.Close()

  efile2, err := NewEntryFile(WithFileInfo(fi2), WithData(fi2Data))
  if err != nil {
    t.Fatal(err)
  }
  os.Remove(file1.Name())
  os.Remove(file2.Name())

  efiles = append(efiles, *efile1)
  efiles = append(efiles, *efile2)
  
  return efiles, err
}


