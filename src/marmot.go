package main
import (
	"github.com/garyburd/redigo/redis"
	"io"
	"crypto/md5"
	"fmt"
)

func hash(in string) string {
	h := md5.New()
	io.WriteString(h,in) 
	var hsh = fmt.Sprintf("%x", h.Sum(nil))
  return hsh
}

func open() redis.Conn {
  r, err := redis.Dial("tcp", ":6379")
  if err != nil {
 	  fmt.Println(err)
  }
  return r 
}

func clean() {
  var r = open()
  r.Do("flushdb")
  r.Flush()
}

func install(site string) {
  fmt.Println("hi")
  var r = open()
  r.Send("hset", "sites", hash(site), site)
  r.Flush()
}

func repo(site string) string {
  var r = open()
	surl, err := redis.String(r.Do("HGET", "sites", hash(site)))
	if err != nil {
		fmt.Println(err)
	}
  fmt.Println(surl)
  return surl
}

func main() {
  clean()
  install("https://github.com/gavinmyers/blog")
  repo("https://github.com/gavinmyers/blog")
/*  var r = open()
  r.Send("SADD", "test_set", "foo")
	r.Send("set","gavin","123") 
	r.Send("SMEMBERS", "test_set")
	r.Flush()
  r.Receive()
  fmt.Println(repo("gavinm.com")) */
}
func GitHandler() {
}

func NilHandler() {
}
