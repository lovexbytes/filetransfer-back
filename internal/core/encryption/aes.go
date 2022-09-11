package encryption

import
(
  "io"
  "errors"
  "strconv"
  "crypto/cipher"
  "crypto/aes"
  "crypto/rand"
)

type aesCipher struct {}

func NewAESCipher() *aesCipher {
  c := &aesCipher{}
  return c
}

func (aesC *aesCipher) Encrypt (key []byte, data []byte) ([]byte, error){
  var ciphertext []byte
  var err error 
   
  block, err := aes.NewCipher(key)
  if err == nil {
    bs := block.BlockSize()
    if len(key) >= bs {
      //Make the cipher text a byte array of size BlockSize + the length of the data
	    ciphertext = make([]byte, bs + len(data))
      //iv is the ciphertext up to the blocksize (16)
	    iv := ciphertext[:bs]
	    _, err = io.ReadFull(rand.Reader, iv)
	    if err == nil {
        //Encrypt the data:
	      stream := cipher.NewCFBEncrypter(block, iv)
	      stream.XORKeyStream(ciphertext[bs:], data)
	    }
    } else {
      err = errors.New("Key is too short for the block size specified. Key length:" + strconv.Itoa(len(key)) + ", block size: " + strconv.Itoa(bs))
    }
  }  
  return ciphertext, err 
}

func (aesC *aesCipher) Decrypt (key []byte, data []byte) ([]byte, error){
  block, err := aes.NewCipher(key)
  if err == nil {
    bs := block.BlockSize()
    if len(data) >= bs{
      iv := data[:bs]
      data = data[bs:]

      //Decrypt the data
      stream := cipher.NewCFBDecrypter(block, iv)
      stream.XORKeyStream(data, data)
    }else{
      err = errors.New("Ciphertext block size is too short")
    }
  }
  return data, err
}
