package cookie

import (
	"net/http"
	"testing"
	"time"
)

var testStore CookieStore = &FileCookieStore{}

func TestFileCookieStore_Init(t *testing.T) {
	testStore.Init("tempCache.txt")
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

	cookies := make([]*http.Cookie, 0)
	cookies = append(cookies, SUV)
	cookies = append(cookies, SNUID)

	testStore.Add(cookies)
}

func TestFileCookieStore_Get(t *testing.T) {
	testStore.Init("tempCache.txt")
	cookies := testStore.Get()
	println(cookies)
}

func TestFileCookieStore_Pop(t *testing.T) {
	testStore.Init("tempCache.txt")
	cookies := testStore.Pop()
	println(cookies)
}

func TestFileCookieStore_Flush(t *testing.T) {
	testStore.Init("tempCache.txt")
	testStore.Flush()
}