package encryption

type Cipher interface{
  Encrypt (key []byte, data []byte) ([]byte, error)
  Decrypt (key []byte, data []byte) ([]byte, error)

}


