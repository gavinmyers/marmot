package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/hoisie/web"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
)

type Payload struct {
	Before  string
	After   string
	Ref     string
	Commits []struct {
		Id        string
		Message   string
		Url       string
		Timestamp string
		Added     []string
		Modified  []string
		Removed   []string
	}
	Repository struct {
		Owner struct {
			Email string
			Name  string
		}
		Description string
		Name        string
		Url         string
	}
}

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

func decode(str string) []byte {
	enc := []byte(str)
	e64 := base64.StdEncoding
	maxDecLen := e64.DecodedLen(len(enc))
	var decBuf = make([]byte, maxDecLen)
	n, _ := e64.Decode(decBuf, enc)
	return decBuf[0:n]
}

func open() redis.Conn {
	r, _ := redis.Dial("tcp", ":6379")
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
	surl, _ := redis.String(r.Do("HGET", "sites", hash(site)))
	return surl
}

func url(repo string, action string, v interface{}) interface{} {
	var r = open()

	client_id, _ := redis.String(r.Do("GET", "client_id"))

	client_secret, _ := redis.String(r.Do("GET", "client_secret"))

	var buffer bytes.Buffer
	buffer.WriteString("https://api.github.com/repos/")
	buffer.WriteString(repo)
	buffer.WriteString("/")
	buffer.WriteString(action)
	buffer.WriteString("?client_id=")
	buffer.WriteString(client_id)
	buffer.WriteString("&client_secret=")
	buffer.WriteString(client_secret)
	res, _ := http.Get(buffer.String())
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return json.Unmarshal(body, &v)
}

func pullFile(repo string, path string) string {
	if path == "" {
		path = "index.html"
	}
	path = "contents/" + path
	var r = open()
	var file GitFile
	var buffer bytes.Buffer
	buffer.WriteString(hash(repo))
	buffer.WriteString(":content")
	url(repo, path, &file)
	r.Send("hset", buffer.String(), hash(path), file.Content)
	r.Flush()
	return file.Content
}

func memFile(repo string, path string) string {
	if path == "" {
		path = "index.html"
	}
	path = "contents/" + path
	var r = open()
	var buffer bytes.Buffer
	buffer.WriteString(hash(repo))
	buffer.WriteString(":content")
	scontent, _ := redis.String(r.Do("hget", buffer.String(), hash(path)))
	r.Flush()
	return scontent
}

func gitFile(repo string, path string) string {
	var scontent = memFile(repo, path)
	if scontent != "" {
		return scontent
	} else {
		return pullFile(repo, path)
	}
	return ""
}

func gitJson(repo string, path string, v interface{}) interface{} {
	return json.Unmarshal(decode(gitFile(repo, path)), &v)
}

//this REALLY shouldn't be here, but mime.TypeByExtension isn't working
func TypeByExtension(file string) string {
	if strings.Contains(file, ".css") {
		return "css"
	} else if strings.Contains(file, ".jpg") {
		return "jpg"
	} else if strings.Contains(file, ".js") {
		return "js"
	} else if strings.Contains(file, ".png") {
		return "png"
	}
	return "text/html"
}

func main() {
	//I shouldn't have to add an extesion, even if I do it still doesn't work
	mime.AddExtensionType("css", "text/css")
	mime.TypeByExtension("test.css")
	clean()
	install("gavinmyers/blog")
	var repo = repository("gavinmyers/blog")
	//get the marmot file
	var config Config
	gitJson(repo, "marmot.json", &config)
	web.Get("/(.*)", func(ctx *web.Context, val string) string {
		ctx.ContentType(TypeByExtension(val))
		return string(decode(gitFile(repo, val)))
	})

	web.Post("/(.*)", func(ctx *web.Context, name string) string {
		var payload Payload
		var repo = repository("gavinmyers/blog")
		json.Unmarshal([]byte(ctx.Params["payload"]), &payload)
		for _, commit := range payload.Commits {
			for _, modified := range commit.Modified {
				pullFile(repo, modified)
			}
			for _, removed := range commit.Removed {
				pullFile(repo, removed)
			}
			for _, added := range commit.Added {
				pullFile(repo, added)
			}
		}
		return ""
	})
	web.Run(config.Url)
}
