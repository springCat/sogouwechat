package cookie

import (
	"net/http"
	"testing"
	"time"
)

var store CookieStore = &FileCookieStore{}

func TestFileCookieStore_Init(t *testing.T) {
	store.Init("tempCache.txt")
	SUV := &http.Cookie{
		Name:    "SUV",
		Value:   "00B68873700A09345D6B9AA23BB99258",
		Path:    "/",
		Domain:  ".sogou.com",
		Expires: time.Now().Add(time.Hour * 24),
	}
	SNUID := &http.Cookie{
		Name:    "SNUID",
		Value:   "92B32DCDB8BD293E62470518B99C41B4",
		Path:    "/",
		Domain:  "weixin.sogou.com",
		Expires: time.Now().Add(time.Hour * 24),
	}

	cookiemap := map[string]*http.Cookie{
		"SUV":SUV,
		"SNUID":SNUID,
	}
	store.Add(cookiemap)
}

func TestFileCookieStore_Get(t *testing.T) {
	store.Init("tempCache.txt")
	cookies := store.Get()
	println(cookies)
}

func TestFileCookieStore_Pop(t *testing.T) {
	store.Init("tempCache.txt")
	cookies := store.Pop()
	println(cookies)
}

func TestFileCookieStore_Flush(t *testing.T) {
	store.Init("tempCache.txt")
	store.Flush()
}