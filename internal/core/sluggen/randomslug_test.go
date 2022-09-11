package generator

import (
	// "context"
	"log"
	"testing"
)

func TestRandomSlugGen(t *testing.T) {
  rsg := NewRandomSlugGen()
  res, err := rsg.GetRandomSlug(0)
  log.Println(res)
  log.Println(err)

  res, err = rsg.GetRandomSlug(3)
  log.Println(res)
  log.Println(err)
  
}
