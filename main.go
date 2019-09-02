package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
)

var cookiePool map[string]*http.Cookie = map[string]*http.Cookie{}
var UA = "Mozilla/6.0 (windows; windows NT) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.132 Safari/537.36"
var vm = otto.New()
type ReqParam struct{
	Key string
	Referer string
	Wxid string
	Tsn string
	Ua string
	Cookies []*http.Cookie
}
func main() {

	//getCookie
	//param := &ReqParam{
	//	Key:"刘备教授",
	//	Ua:UA,
	//}
	//response, err:= getCookie("刘备教授", nil)
	//assertOk(err)
	//
	//cookies := response.Cookies()
	//for _, cookie := range cookies {
	//	cookiePool[cookie.Name]=cookie;
	//}

	param := &ReqParam{
		Key:"刘备教授",
		Wxid:"oIWsFtx2SU5am12hfw0hb6qYgUXg",
		Tsn:"1",
		Ua:UA,
		Cookies:buildCookies(),
	}
	resp, err := search(param)
	assertOk(err)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	assertOk(err)

	val, exists := doc.Find(".news-list li .txt-box h3 a").Attr("href")
	if(exists){

		param.Key = genContentReqUrl(val)
		param.Referer = "https://weixin.sogou.com/weixin?type=2&ie=utf8&query=刘备教授&tsn=1&wxid=oIWsFtx2SU5am12hfw0hb6qYgUXg"

		resp, err = queryContentUrl(param)
		assertOk(err)

		contentUrl := parseContentUrl(resp)
		log.Println(contentUrl)

		resp, err := sogouWechatGet(contentUrl, param.Ua, param.Referer, nil)
		assertOk(err)

		bytes, err := ioutil.ReadAll(resp.Body)
		assertOk(err)

		doc, err = goquery.NewDocumentFromReader(resp.Body)

		ioutil.WriteFile("/Users/springcat/Desktop/temp.html",bytes,os.ModePerm)
		log.Println(string(bytes))
	}

}
/**
 1 SUV
 2 SNUID
 */
func buildCookies() []*http.Cookie {
	SNUID := &http.Cookie{
		Name:"SUV",
		Value:"00B68873700A09345D6B9AA23BB99258",
		Path:"/",
		Domain:".sogou.com",
		Expires:time.Now().Add(time.Hour*24),
	}

	SUID := &http.Cookie{
		Name:"SNUID",
		Value:"92B32DCDB8BD293E62470518B99C41B4",
		Path:"/",
		Domain:"weixin.sogou.com",
		Expires:time.Now().Add(time.Hour*24),
	}

	cookies := make([]*http.Cookie, 0)
	cookies = append(cookies, SNUID)
	cookies = append(cookies, SUID)
	return cookies

}

func search(param *ReqParam) (resp *http.Response, err error){
	url := "https://weixin.sogou.com/weixin?type=2&ie=utf8&query="+param.Key+"&tsn="+param.Tsn+"&wxid="+param.Wxid
	return sogouWechatGet(url, param.Ua,url, param.Cookies)
}

func queryContentUrl(param *ReqParam) (resp *http.Response, err error){
	genContentReqUrl(param.Key)
	url := "https://weixin.sogou.com"+param.Key
	log.Println(url)
	return sogouWechatGet(url, param.Ua,param.Referer, param.Cookies)
}

/**
   function() {
                var b = Math.floor(100 * Math.random()) + 1
                  , a = this.href.indexOf("url=")
                  , c = this.href.indexOf("&k=");
                -1 !== a && -1 === c && (a = this.href.substr(a + 4 + parseInt("21") + b, 1),
                this.href += "&k=" + b + "&h=" + a)
            }
 */
func genContentReqUrl(rawUrl string) string{
	script :=`
		var href = "`+rawUrl+`";
		var b = Math.floor(100 * Math.random()) + 1
		, a = href.indexOf("url=")
		, c = href.indexOf("&k=");
		-1 !== a && -1 === c && (a = href.substr(a + 4 + parseInt("21") + b, 1),
			href += "&k=" + b + "&h=" + a)
		
	`
	value, err := vm.Run(script)
	assertOk(err)

	log.Println("genContentReqUrl value:"+value.String())
	return value.String()
}

func parseContentUrl(resp *http.Response) string {


	doc, err := goquery.NewDocumentFromReader(resp.Body)
	assertOk(err)

	log.Println("parseContentUrl doc:"+doc.Text())

	script := doc.Find("script").Text()

	log.Println("parseContentUrl script:"+script)

	script = strings.Replace(script,"window.location.replace(url)","",-1)

	value, err := vm.Run(script)
	assertOk(err)

	log.Println("genContentReqUrl url:"+value.String())
	return value.String()
}


func getCookie(param ReqParam) (resp *http.Response, err error){
	url := "https://weixin.sogou.com/weixin?type=2&query="+param.Key
	return sogouWechatGet(url, param.Ua,url, param.Cookies)
}
//
//func retryByChangeCookie(url string,ua string,referer string,cookies []*http.Cookie)  (resp *http.Response, err error){
//	request, err := http.NewRequest("GET", url, nil)
//	assertOk(err)
//
//	request.Header.Add("User-Agent", ua)
//	request.Header.Add("Referer", referer)
//	if(cookies != nil) {
//		for _, cookie := range cookies {
//			request.AddCookie(cookie)
//		}
//	}
//	return http.DefaultClient.Do(request)
//}

func sogouWechatGet(url string,ua string,referer string,cookies []*http.Cookie)  (resp *http.Response, err error){
	request, err := http.NewRequest("GET", url, nil)
	assertOk(err)

	request.Header.Add("User-Agent", ua)
	request.Header.Add("Referer", referer)
	if(cookies != nil) {
		for _, cookie := range cookies {
			request.AddCookie(cookie)
		}
	}
	return http.DefaultClient.Do(request)
}


func assertOk(err error) {
	if err != nil {
		log.Panic(err)
	}
}
