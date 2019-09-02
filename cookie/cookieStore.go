package cookie

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

/**
	FIFO testStore
 */
type CookieStore interface {
	Get() []*http.Cookie
	Pop() []*http.Cookie
	Add([]*http.Cookie) error
	Flush()
	Init(cacheName string) error

}

type FileCookieStore struct {
	FilePath   string
	isInit     bool
	queue      [][]*http.Cookie
	len        int
}

func (store *FileCookieStore) Init(cacheName string) error {
	store.FilePath = cacheName
	bytes, e := ioutil.ReadFile(store.FilePath)
	if e == nil {
		json.Unmarshal(bytes, &store.queue)
	}
	store.len = len(store.queue)
	store.isInit = true
	return e;
}

func (store *FileCookieStore) Get() []*http.Cookie{
	store.checkInit()
	if !store.isInit{
		panic("FileCookieStore not init")
	}
	if(store.len > 0){
		return store.queue[0]
	}
	return nil
}

func (store *FileCookieStore) Pop() []*http.Cookie{
	store.checkInit()
	if(store.len > 0){
		last := store.queue[0]
		store.queue = store.queue[1:store.len]
		store.len = store.len - 1
		bytes, e := json.Marshal(store.queue)
		if(e == nil) {
			ioutil.WriteFile(store.FilePath,bytes,os.ModePerm)
		}
		return last
	}
	return nil
}

func (store *FileCookieStore) Add(cookies []*http.Cookie) error{
	store.checkInit()
	if(store.len == 0) {
		store.queue = make([][]*http.Cookie,0)
	}
	store.queue = append(store.queue,cookies)
	store.len = store.len + 1
	bytes, e := json.Marshal(store.queue)
	if(e == nil) {
		ioutil.WriteFile(store.FilePath,bytes,os.ModePerm)
	}
	return e
}

func (store *FileCookieStore) Flush(){
	store.checkInit()
	store.queue = make([][]*http.Cookie,0)
	store.len = 0
	ioutil.WriteFile(store.FilePath,nil,os.ModePerm)
}

func (store *FileCookieStore) checkInit(){
	if !store.isInit{
		panic("FileCookieStore not init")
	}
}