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

func dial() redis.Conn {
  r, err := redis.Dial("tcp", ":6379")
  if err != nil {
 	  fmt.Println(err)
  }
  return r 
}
func get(r redis.Conn, k string) string {
	reply, err := redis.String(r.Do("GET", k))
  if err != nil {
 	  fmt.Println(err)
  }
  return reply
}

func clean() {
  var r = dial()
  r.Do("flushdb")
  r.Flush()
}

func install(site string) {
  fmt.Println("hi")
  var r = dial()
  r.Send("hset", "sites", hash(site), site)
  r.Flush()
}

func repo(site string) string {
  var r = dial()
	surl, err := redis.String(r.Do("HGET", "sites", hash(site)))
	if err != nil {
		fmt.Println(err)
	}
  fmt.Println(surl)
  return surl
}

func main() {
  clean()
  install("gavinm.com")
  repo("gavinm.com")
  var r = dial()
  r.Send("SADD", "test_set", "foo")
	r.Send("set","gavin","123") 
	r.Send("SMEMBERS", "test_set")
	r.Flush()
  r.Receive()
  fmt.Println(repo("gavinm.com"))
}
func GitHandler() {
}

func NilHandler() {
}
