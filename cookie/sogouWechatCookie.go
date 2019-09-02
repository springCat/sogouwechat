package cookie

import (
	"log"
	"net/http"
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

