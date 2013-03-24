package main

import (
  "encoding/json"
  "github.com/garyburd/redigo/redis"
	"net/http"
	"fmt"
	"io/ioutil"
	"log"
)

type Payload struct {
  Forced bool
  Repository struct {
    Private bool
    Owner struct {
      Email string
      Name string
      Has_downloads bool
      Stargazers int
      Id string
      Watchers int
      Master_branch string
      Has_wiki bool
      Description string
      Fork bool
    }
  }
}

func main() {
  c, err := redis.Dial("tcp", ":6379")
  defer c.Close()
  if err != nil {
    fmt.Println(err)
  }
  c.Send("SADD", "myset", "{hello:0}")
  c.Send("SADD", "myset", "{world:1}")
  c.Send("SMEMBERS", "myset")
  c.Flush()
  c.Receive()
  //c.Receive()
	http.HandleFunc("/git", GitHandler)
	http.HandleFunc("/favicon.ico", NilHandler)
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe("unstable.gavinm.com:8080", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
  c, err := redis.Dial("tcp", ":6379")
  defer c.Close()
  if err != nil {
    fmt.Println(err)
  }
	resp, err := http.Get("https://raw.github.com/gavinmyers/resume/master/README.md")
	if err != nil {
		fmt.Fprint(w,err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Fprint(w,string(body))
  err = c.Send("SMEMBERS", "myset")
  err = c.Send("get", "payload")
  if err != nil {
    fmt.Println(err)
  }
  c.Flush()
  // both give the same return value!?!?
  // reply, err := c.Receive()
  reply, err := redis.MultiBulk(c.Receive())
  if err != nil {
    fmt.Println(err)
  }
  for _, x := range reply {
    var v, ok = x.([]byte)
    if ok {
      fmt.Println(string(v))
    }
  }
  //fmt.Printf("%#v\n", reply)
}

func GitHandler(w http.ResponseWriter, r *http.Request) {
  c, err := redis.Dial("tcp", ":6379")
  defer c.Close()
  if err != nil {
    fmt.Println(err)
  }
  payload := r.FormValue("payload")
  c.Send("set", "payload", payload)
  c.Flush()
  c.Receive()
	fmt.Fprint(w, "This is Git handler! " + payload)
}

func NilHandler(w http.ResponseWriter, r *http.Request) {
}

