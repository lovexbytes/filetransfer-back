package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"filetransfer-back/internal/core/archive"
	"filetransfer-back/internal/core/encryption"
	"filetransfer-back/internal/config"
	"github.com/google/uuid"
)

func TestSecuredStorage(t *testing.T){
    /**
    1. create ss
    2. archive ss
    3. encrypt ss
    4. check content
    5. decrypt ss
    6. unarchive ss
    7. check content
  */

	//init test EntryFile array
	efiles, err := initTestEfiles(t)
	if err != nil {
		t.Fatal(err)
	}

	//create archiver
	archiver := archive.NewZipArchiver()
	
	//create cipher
	cipher := encryption.NewAESCipher()
	
	//create ss
	ss, err := NewSecuredStorage(
															WithSSFileEntries(efiles),
															WithSSArchiver(archiver),
															WithSSCipher(cipher),
															)
	if err != nil {
		t.Fatal(err)
	}
	
	for _, efile := range(ss.EntryFiles()) {
		if len(efile.Data()) == 0 {
      t.Fatal(errors.New("efile " + efile.FileInfo().Name() + " data is empty"))
    }
  }

	//archive ss
	err = ss.Archive()
	if err != nil {
		t.Fatal(err)
	}

	//encrypt ss
	err = ss.Encrypt()
	if err != nil {
		t.Fatal(err)
	}

	//check that ss files are no longer available
	if len(ss.EntryFiles()) != 0 {
		t.Fatal("ss.EntryFiles() must be empty after archive and encryption proccesses")
	}

	//check that data buffer is not empty as it contains data that has to be saved
	bytedata, err := ss.ByteData()
	if err != nil {
		t.Fatal(err)
	}

	if len(bytedata) == 0 {
		t.Fatal("ss.bytedata (encrypted and archived data) is empty")
	}
	
	//creating file to store archived & encrypted data
	ssId := ss.Id() //storing it for unpacking
	containerFilename := ssId + "." + config.SecureStorageExtension()
	container, err := os.Create(containerFilename)
	if err != nil {
		t.Fatal(err)
	}
	
	encryptedByteData, err := ss.ByteData()
	if err != nil {
		t.Fatal(err)
	}

	n, err := container.Write(encryptedByteData)
	if n == 0 {
		t.Fatal("0 bytes written during test save of archived & encrypted SS")
	}

	err = container.Close()
	if err != nil {
		t.Fatal(err)
	}

	//unpacking ss
	container, err = os.OpenFile(ssId + ".ss", os.O_RDONLY, 0755) //opening with existing id
	if err != nil {
		t.Fatal(err)
	}

	unpackEncryptedByteData, err := io.ReadAll(container) 
	if err != nil {
		t.Fatal(err)
	}

	//removing encrypted ss
	err = os.Remove(container.Name())
	if err != nil {
		t.Fatal(err)
	}

	ssUUID, err := uuid.ParseBytes([]byte(ssId))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(len(unpackEncryptedByteData))
	unpackSS, err := NewSecuredStorage(WithSSUUID(ssUUID), WithByteData(unpackEncryptedByteData), WithSSArchiver(archiver), WithSSCipher(cipher))
	if err != nil {
		t.Fatal(err)
	}

	err = unpackSS.Unpack()
	if err != nil {
		t.Fatal(err)
	}

	//creating a file to store unpacked ss data
	file, err := os.Create(unpackSS.Id() + ".zip")
	if err != nil {
		t.Fatal(err)
	}

	readyData := unpackSS.bytedata
	
	n, err = file.Write(readyData)
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Fatal("0 bytes were written to the unpacked ss file")
	}

	file.Close()
	

}


//util
func initTestEfiles (t *testing.T) ([]archive.EntryFile, error) {
  efiles := []archive.EntryFile{}

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

  efile1, err := archive.NewEntryFile(archive.WithFileInfo(fi1), archive.WithData(fi1Data))
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

  efile2, err := archive.NewEntryFile(archive.WithFileInfo(fi2), archive.WithData(fi2Data))
  if err != nil {
    t.Fatal(err)
  }
  os.Remove(file1.Name())
  os.Remove(file2.Name())

  efiles = append(efiles, *efile1)
  efiles = append(efiles, *efile2)
  
  return efiles, err
}


