package simpleapm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Agent struct {
	APMUrl    string
	AppId     string
	AppSecret string
}

func (self *Agent) print(message string) {
	log.Println(message)
}

func (self *Agent) postHTTP(post *PostHttp) {

	// generate json string
	b, err := json.Marshal(map[string]interface{}{"post": post})
	if err != nil {
		fmt.Println(err)
		return
	}

	// do request in seperated thread
	go self.doRequest(string(b))
}

func (self *Agent) doRequest(data string) {

	url := strings.TrimRight(self.APMUrl, "/") + "/posts"

	var jsonBytes = []byte(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-apm-app-secret", self.AppSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))
}

func (self *Agent) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// customize the response object
		resp := NewCustomResponseWriter(w)

		// start measure
		now := time.Now()

		// do next in the chain
		next.ServeHTTP(resp, r)

		// finished request
		timeBench := float64(time.Now().Sub(now).Nanoseconds()) / float64(1e6)

		// create post struct
		post := &PostHttp{Post: &Post{
			Type:      "http",
			Status:    resp.StatusCode(),
			StartTime: now,
			EndTime:   time.Now()},
			Data: PostData{
				Path:     r.URL.Path,
				Method:   r.Method,
				Duration: timeBench,
				Query:    r.URL.RawQuery,
				HttpCode: resp.StatusCode()}}

		self.postHTTP(post)
		//log.Printf("%.2fms\t| %d | %s: %s\n", timeBench, resp.StatusCode(), r.Method, r.URL)
	})
}
