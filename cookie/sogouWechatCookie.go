package cookie

import (
	"log"
	"net/http"
	"time"

	"org.springcat/sougoWeixin/tools"
)

var store CookieStore

func init()  {
	store = &FileCookieStore{}
	store.Init("cookiestore")
}
//--------------------------------- cookie handle end---------------------------------
func FetchCookie(ua string){
	url := "https://weixin.sogou.com/weixin?type=2&query=springcat"
	resp, err := tools.SogouWechatGet(url, ua, url, nil)

	if(err != nil){
		log.Println("fetchCookie:"+err.Error())
	}

	cookies := make([]*http.Cookie, 0)
	for _,v := range resp.Cookies(){
		if v.Name == "SUV" || v.Name == "SNUID" {
			cookies = append(cookies,v)
		}
	}

	SUV := &http.Cookie{
		Name:    "SUV",
		Value:   "00B68873700A09345D6B9AA23BB99258",
		Path:    "/",
		Domain:  ".sogou.com",
		Expires: time.Now().Add(time.Hour * 24),
	}
	SNUID := &http.Cookie{
		Name:    "SNUID",
		Value:   "0626BE592B29BDB63659830F2C4AD60B",
		Path:    "/",
		Domain:  ".sogou.com",
		Expires: time.Now().Add(time.Hour * 24),
	}
	cookies = append(cookies,SUV)
	cookies = append(cookies,SNUID)

	if len(cookies) > 0{
		store.Add(cookies)
	}
}

func RetryByChangeCookie(url string, ua string, referer string, cookies []*http.Cookie) (resp *http.Response, err error) {
	response, e := tools.SogouWechatGet(url, ua, referer, cookies)
	tools.AssertOk(e)
	if (response.StatusCode == 301) {
		log.Println("retryByChangeCookie change cookie url:" + url)
		store.Pop()
		cookies = store.Get()
	}
	return RetryByChangeCookie(url, ua, referer, cookies)
}

func GetCookie() []*http.Cookie {
	return store.Get()
}

