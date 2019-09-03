package tools

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
	resp,err = http.DefaultClient.Do(request)

	AssertOk(err)
	htttlog(request,resp)

	return
}

func htttlog(request *http.Request,response *http.Response)  {
	log.Println("SogouWechatGet request url:"+request.URL.String())
	log.Println("SogouWechatGet request ua:"+request.UserAgent())
	log.Println("SogouWechatGet request referer:"+request.Header.Get("referer"))
	bytes, e := json.Marshal(request.Cookies())
	AssertOk(e)
	log.Println("SogouWechatGet request cookie:"+string(bytes))
	log.Println("SogouWechatGet response status:"+response.Status)
	log.Println("SogouWechatGet response content-length:"+strconv.FormatInt(response.ContentLength,10))
	bytes, e = json.Marshal(response.Cookies())
	AssertOk(e)
	log.Println("SogouWechatGet response cooke:"+string(bytes))
}

func AssertOk(err error) {
	if err != nil {
		log.Panic(err)
	}
}

