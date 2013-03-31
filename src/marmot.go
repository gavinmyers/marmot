package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io"
	"io/ioutil"
	"net/http"
)

type GitFile struct {
	Sha      string
	Name     string
	Path     string
	Type     string
	Url      string
	Git_url  string
	Html_url string
	Content  string
	Encoding string
}

type Config struct {
	Description string
	Url         string
}

func hash(in string) string {
	h := md5.New()
	io.WriteString(h, in)
	var hsh = fmt.Sprintf("%x", h.Sum(nil))
	return hsh
}

func decode(str string, v interface{}) interface{} {
	enc := []byte(str)
	e64 := base64.StdEncoding
	maxDecLen := e64.DecodedLen(len(enc))
	var decBuf = make([]byte, maxDecLen)
	n, err := e64.Decode(decBuf, enc)
	_ = err
	return json.Unmarshal(decBuf[0:n], &v)
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
	//r.Do("flushdb")
	r.Flush()
}

func install(site string) {
	var r = open()
	r.Send("hset", "sites", hash(site), site)
	r.Flush()
}

func repository(site string) string {
	var r = open()
	surl, err := redis.String(r.Do("HGET", "sites", hash(site)))
	if err != nil {
		fmt.Println("repository")
		fmt.Println(err)
	}
	return surl
}

func url(repo string, action string, v interface{}) interface{} {
	//https://api.github.com/repos/gavinmyers/blog/contents/
	var buffer bytes.Buffer
	buffer.WriteString("https://api.github.com/repos/")
	buffer.WriteString(repo)
	buffer.WriteString("/")
	buffer.WriteString(action)
	res, err := http.Get(buffer.String())
	if err != nil {
		fmt.Println("url")
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("url")
		fmt.Println(err)
	}
	return json.Unmarshal(body, &v)
}

func gitFile(repo string, path string, v interface{}) interface{} {
	var r = open()
	var buffer bytes.Buffer
	buffer.WriteString(hash(repo))
	buffer.WriteString(":content")
	scontent, err := redis.String(r.Do("hget", buffer.String(), hash(path)))
	r.Flush()
	if err != nil {
		fmt.Println(err)
	}
	if scontent != "" {
		return decode(scontent, &v)
	} else {
		var file GitFile
		url(repo, path, &file)
		r = open()
		r.Send("hset", buffer.String(), hash(path), file.Content)
		r.Flush()
		return decode(file.Content, &v)
	}
	return nil
}

func main() {
	clean()
	install("gavinmyers/blog")
	var repo = repository("gavinmyers/blog")
	//get the marmot file
	var config Config
	gitFile(repo, "contents/marmot.json", &config)
	fmt.Println(config.Url)
	http.HandleFunc("/favicon.ico", NilHandler)
	http.ListenAndServe(config.Url, nil)
	/*  var r = open()
	  r.Send("SADD", "test_set", "foo")
		r.Send("set","gavin","123") 
		r.Send("SMEMBERS", "test_set")
		r.Flush()
	  r.Receive()
	  fmt.Println(repo("gavinm.com")) */
}
func GitHandler(w http.ResponseWriter, r *http.Request) {
}

func NilHandler(w http.ResponseWriter, r *http.Request) {
}
