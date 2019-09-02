package cookie

import (
	"log"
	"net/http"
	"time"

	"org.springcat/sougoWeixin/tools"
)

//--------------------------------- cookie handle end---------------------------------
func getCookie(param tools.ReqParam) (resp *http.Response, err error) {
	url := "https://weixin.sogou.com/weixin?type=2&query=" + param.Key
	return tools.SogouWechatGet(url, param.Ua, param.Referer, param.Cookies)
}

/**
 1 SUV
 2 SNUID
 */
func BuildCookies() []*http.Cookie {
	SNUID := &http.Cookie{
		Name:    "SUV",
		Value:   "00B68873700A09345D6B9AA23BB99258",
		Path:    "/",
		Domain:  ".sogou.com",
		Expires: time.Now().Add(time.Hour * 24),
	}

	SUID := &http.Cookie{
		Name:    "SNUID",
		Value:   "92B32DCDB8BD293E62470518B99C41B4",
		Path:    "/",
		Domain:  "weixin.sogou.com",
		Expires: time.Now().Add(time.Hour * 24),
	}

	cookies := make([]*http.Cookie, 0)
	cookies = append(cookies, SNUID)
	cookies = append(cookies, SUID)
	return cookies
}

func NewCookies() []*http.Cookie {
	SNUID := &http.Cookie{
		Name:    "SUV",
		Value:   "00B68873700A09345D6B9AA23BB99258",
		Path:    "/",
		Domain:  ".sogou.com",
		Expires: time.Now().Add(time.Hour * 24),
	}

	SUID := &http.Cookie{
		Name:    "SNUID",
		Value:   "92B32DCDB8BD293E62470518B99C41B4",
		Path:    "/",
		Domain:  "weixin.sogou.com",
		Expires: time.Now().Add(time.Hour * 24),
	}

	cookies := make([]*http.Cookie, 0)
	cookies = append(cookies, SNUID)
	cookies = append(cookies, SUID)
	return cookies
}

func RetryByChangeCookie(url string, ua string, referer string, cookies []*http.Cookie) (resp *http.Response, err error) {
	response, e := tools.SogouWechatGet(url, ua, referer, cookies)
	tools.AssertOk(e)
	if (response.StatusCode == 301) {
		log.Println("retryByChangeCookie change cookie url:" + url)
		cookies = NewCookies();
	}
	return RetryByChangeCookie(url, ua, referer, cookies)
}

//--------------------------------- cookie handle end---------------------------------

