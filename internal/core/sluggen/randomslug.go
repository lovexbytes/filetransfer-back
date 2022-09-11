package generator

import (
  "os"
  "errors"
  "strings"
  "bufio"
  "io"
  "context"
  "math/rand"
  "time"
  "encoding/json"
  "filetransfer-back/internal/core/cache"
  "filetransfer-back/internal/config"
)

type randomSlugGen struct {
  ctx context.Context
}

func NewRandomSlugGen () *randomSlugGen {
  rsg := &randomSlugGen{
    ctx: context.TODO(),
  }
  return rsg
}

func (rsg *randomSlugGen) GetRandomSlug (wordsNum int) (string, error) {
  var err error
  var dict []string
  var slug string
  
  if wordsNum > 0 {
    //random seeding to generate new values for each run 
	  rand.Seed(time.Now().UnixNano())
    var words []string
    var sep string = "-"

    dict, err = getDictionary(rsg.ctx)
    if err == nil {
      linesNum := len(dict)
      if linesNum > 0 {
        var bufLine string
        for i := 0; i < wordsNum; i++ {  
          bufLine = dict[rand.Intn(linesNum) - 1]
          words = append(words, bufLine)
        }
        slug = strings.ToLower(strings.Join(words, sep))
      }
    }  
  } else {
    err = errors.New("Number of words cannot be zero")
  }
  return slug, err
}

func getDictPath () string {
  sep := string(os.PathSeparator)
  dictPath := strings.Join([]string{"", "usr", "share", "dict", "words"}, sep) 
  return dictPath
}

func getDictionary (ctx context.Context) ([]string, error){
  var err error
  var dict []string
  strDict, err := cache.RedisInstance(ctx).Get(ctx, config.RedisDictionaryKey()).Result()
  if err != nil {
    dict, err = getDictFileContent()
    if err == nil {
      err = cacheDict(ctx, dict)
    }
  }else{
    err = json.Unmarshal([]byte(strDict), &dict)
  }
  return dict, err
}

func getDictFileContent () ([]string, error) {
  var err error
  var dictArr []string

  dictPath := getDictPath()
  dict, err := os.Open(dictPath)
  reader := bufio.NewReader(dict)
  for {
    line, _, readErr := reader.ReadLine() 
    if readErr != nil {
      if readErr != io.EOF{
        err = readErr
      }
      break
    }else{
      dictArr = append(dictArr, string(line))
    }
  }
  return dictArr, err
}

func cacheDict(ctx context.Context, dictArr []string) error{
  var err error
  //serialize dictionary array and put it to redis
  if err == nil {
    var serializedDict []byte 
    serializedDict, err = json.Marshal(dictArr)
    if err == nil {
      err = cache.RedisInstance(ctx).Set(ctx, config.RedisDictionaryKey(), serializedDict, time.Hour).Err()
    }
  }
  return err
}

