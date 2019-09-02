package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
	"org.springcat/sougoWeixin/tools"
	"org.springcat/sougoWeixin/cookie"
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

	cookie.FetchCookie(UA)

	param := &ReqParam{
		Key:"刘备教授",
		Wxid:"oIWsFtx2SU5am12hfw0hb6qYgUXg",
		Tsn:"3",
		Ua:UA,
		Referer :"https://weixin.sogou.com/weixin?type=2&ie=utf8&query=刘备教授&tsn=1&wxid=oIWsFtx2SU5am12hfw0hb6qYgUXg",
		Cookies:cookie.GetCookie(),
	}
	resp, err := search(param)
	tools.AssertOk(err)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	tools.AssertOk(err)

	val, exists := doc.Find(".news-list li .txt-box h3 a").Attr("href")
	if(exists){
		param.Key = val
		contentUrl := QueryContentUrl(param)
		param.Key = contentUrl
		getContent(param)
	}
}



/**
 * 1st search the key
 */
func search(param *ReqParam) (resp *http.Response, err error){
	url := "https://weixin.sogou.com/weixin?type=2&ie=utf8&query="+param.Key+"&tsn="+param.Tsn+"&wxid="+param.Wxid
	log.Println("search url:"+url)
	return tools.SogouWechatGet(url, param.Ua,param.Referer, param.Cookies)
}

/**
 * ------------------------------------------------------------------------------
 * get the contentUrl
 */

func QueryContentUrl(param *ReqParam) string {
	reqUrl := genContentReqUrl(param.Key)
	url := "https://weixin.sogou.com"+reqUrl
	log.Println("queryContentUrl url:"+url)
	resp, err := tools.SogouWechatGet(url, param.Ua, param.Referer, param.Cookies)
	tools.AssertOk(err)
	contentUrl := parseContentUrl(resp)
	return contentUrl
}

/**
 * handle method js in page
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
	tools.AssertOk(err)
	log.Println("genContentReqUrl value:"+value.String())
	return value.String()
}

func parseContentUrl(resp *http.Response) string {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	tools.AssertOk(err)
	log.Println("parseContentUrl doc:"+doc.Text())
	script := doc.Find("script").Text()
	log.Println("parseContentUrl script:"+script)
	script = strings.Replace(script,"window.location.replace(url)","",-1)
	value, err := vm.Run(script)
	tools.AssertOk(err)
	log.Println("genContentReqUrl url:"+value.String())
	return value.String()
}

/**
 * ------------------------------------------------------------------------------
 * get the contentUrl
 */

/**
 * 3rd get content
 */

func getContent(param *ReqParam)  {
	resp, err := tools.SogouWechatGet(param.Key, param.Ua, param.Referer, nil)
	tools.AssertOk(err)
	bytes, err := ioutil.ReadAll(resp.Body)
	tools.AssertOk(err)
	ioutil.WriteFile("/Users/springcat/Desktop/temp.html",bytes,os.ModePerm)
}
