package encryption

import (
  "testing"
  "os"
  "io/ioutil"
  "fmt"
  "strconv"
)

func TestDefaultAesBehaviour(t *testing.T) {
  //create default aes instance
  aes := NewAESCipher()
  
  
  //create test files and fill it with content
	testFilenames := []string{"testfile1.txt", "testfile2.txt"}
  var files []*os.File
	for i, filename := range testFilenames {
		err := os.WriteFile(filename, []byte("test" + strconv.Itoa(i)), 0600)
		if err != nil {
			t.Fatal(err)
		}
    f, err := os.OpenFile(filename, os.O_RDWR, 0600)
    if err != nil {
      t.Fatal(err)
    }
		files = append(files, f)
		defer f.Close()
	} 
  
  //encrypt files
  testKey := []byte("jo8aegPKlvAAEa2M3l5sePv8JJMnqlTD")
  encryptedFiles := [][]byte{}

  for _, file := range files {
    fileData, err := ioutil.ReadFile(file.Name())
    if err != nil {
      t.Fatal(err)
    }
    encryptedFileData, err := aes.Encrypt(testKey, fileData)
    if err != nil {
      t.Fatal(err)
    }
    if string(fileData) == string(encryptedFileData) {
      t.Fatal("Unencrypted file data equals to the encrypted one")
    }

    decryptedFileData, err := aes.Decrypt(testKey, encryptedFileData)
    if err != nil {
      t.Fatal(err)
    }

    fmt.Println("source=", string(fileData), " | encrypted=", string(encryptedFileData), " | decrypted=",string(decryptedFileData))
    encryptedFiles = append(encryptedFiles, encryptedFileData)
  }

  //replace existing files with encrypted ones
  if len(testFilenames) > len(encryptedFiles){
    t.Fatal("Number of filenames exceeds number of encryptedFiles")
  }

  for i, filename := range testFilenames {
    encryptedFileData := encryptedFiles[i]
    os.WriteFile(filename, encryptedFileData, 600)
  }

   //delete files
	for _, file := range files {
		err := os.Remove(file.Name())
		if err != nil {
			t.Fatal(err)
		}
	}

}


