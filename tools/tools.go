package tools

import (
	"log"
	"net/http"
)

type ReqParam struct {
	Key     string
	Referer string
	Wxid    string
	Tsn     string
	Ua      string
	Cookies []*http.Cookie
}


func SogouWechatGet(url string,ua string,referer string,cookies []*http.Cookie)  (resp *http.Response, err error){
	request, err := http.NewRequest("GET", url, nil)
	AssertOk(err)
	request.Header.Add("User-Agent", ua)
	request.Header.Add("Referer", referer)
	if(cookies != nil) {
		for _, cookie := range cookies {
			request.AddCookie(cookie)
		}
	}
	return http.DefaultClient.Do(request)
}

func AssertOk(err error) {
	if err != nil {
		log.Panic(err)
	}
}

