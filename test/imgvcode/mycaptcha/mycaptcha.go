package main

import (
	"net/http"
	"github.com/dchest/captcha"
	"time"
	"fmt"
)


func pic(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "image/png")
	id := captcha.NewLen(6)
	d := []byte{byte(1),byte(2),byte(3),byte(4),byte(5)}
	_, err := captcha.NewImage(id, d, 200, 80).WriteTo(w)
	if nil != err{
		for _, b := range d{
			fmt.Print(int(b))
		}
	}
	//captcha.WriteImage(w,id, 260, 80)
}
func index(w http.ResponseWriter, req *http.Request) {
	str := "<meta charset=\"utf-8\"><h3>golang 图片验证码例子</h3><img border=\"1\" src=\"/pic\" alt=\"图片验证码\" onclick=\"this.src='/pic'\" />"
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(str))
}
func main() {
	//http.Handle("/pic", captcha.Server(captcha.StdWidth, captcha.StdHeight))
	http.HandleFunc("/pic", pic)
	http.HandleFunc("/", index)
	s := &http.Server{
		Addr:           ":8081",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}