package config

const (
  //redis connection
  _redisIp = "127.0.0.1"
  _redisPort = "6379"
  _redisPass = ""
  _redisDbi = 0
  //redis key values
  _redisDictKey = "DICTIONARY"
)

//redis connection
func RedisIp() string {
  return _redisIp
}

func RedisPort() string {
  return _redisPort
}

func RedisPass() string {
  return _redisPass
}

func RedisDbi() int {
  return _redisDbi
}


//redis key values
func RedisDictionaryKey() string {
  return _redisDictKey
}
