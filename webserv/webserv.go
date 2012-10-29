package main

import (
	//"os"
	"fmt"
	"log"
	"time"
	"io/ioutil"
	"net/http"
)

type Page struct {
	Path           string
	Body           []byte
	RemainingUsage int
	Creation       time.Time
}

const (
	CacheMaxAge       = 60 // sec
	CacheMaxUsage     = 10 // count down
	CachePeriodicTidy = 1  // minutes
)

var (
	pageCache map[string]*Page
	rootPath  string           = "../negentropia"
)

func cacheTidyPeriodic() {
	for {
		//fmt.Printf("tiding up cache\n")
		log.Printf("tiding up cache\n")

		// mutex write lock
		// tidy up cache
		// mutex write unlock

		time.Sleep(CachePeriodicTidy * time.Minute)
	}
}

func absPath(path string) string {
	return rootPath + path
}

func loadPage(path string) (*Page, error) {

	fullPath := absPath(path)
	body, err := ioutil.ReadFile(fullPath)
	if err != nil {
		//fmt.Printf("loadPage: fullPath=%s not found\n", fullPath)
		log.Printf("loadPage: fullPath=%s not found\n", fullPath)
		return nil, err
	}

	return &Page{Path: path, Body: body, RemainingUsage: CacheMaxUsage, Creation: time.Now()}, nil
}

func cacheAge(page *Page) time.Duration {
	return time.Now().Sub(page.Creation) / time.Second // Duration in seconds
}

func pageNotFound(path string) *Page {
	//fmt.Printf("path=%s not found\n", path)
	log.Printf("path=%s not found\n", path)
	return &Page{Path: path, Body: []byte("File not found: [" + path + "]"), Creation: time.Now()}
}

func handler(w http.ResponseWriter, r *http.Request) {

	var (
		page *Page
		err  error
		ok   bool
	)

	//path := r.URL.Path[1:]
	path := r.URL.Path

	//fmt.Printf("handler path=%s\n", path)
	log.Printf("handler path=%s\n", path)	
	
	/*
	// root url ?
	if path[:1] == "/" && len(path[1:]) == 0 {
		path = "/index.html"
	}
	*/
	
	fullPath := absPath(path)
	
	//fmt.Printf("handler url=%s fullPath=%s\n", path, fullPath)
	log.Printf("handler url=%s fullPath=%s\n", path, fullPath)

	/*
	f, err := os.Open(fullPath)
	var modtime time.Time
	http.ServeContent(w, r, fullPath, modtime, f)
	*/
	http.ServeFile(w, r, fullPath)	

	var delay time.Duration = 20
	log.Printf("handler url=%s fullPath=%s sleeping %d secs", path, fullPath, delay)
	time.Sleep(delay * time.Second)
	
	return

	page, ok = pageCache[path]
	if ok {
		age := cacheAge(page)
		page.RemainingUsage--
		//fmt.Printf("handler path=%s cache HIT countdown=%d age=%d/%d\n", path, page.RemainingUsage, age, CacheMaxAge)
		log.Printf("handler path=%s cache HIT countdown=%d age=%d/%d\n", path, page.RemainingUsage, age, CacheMaxAge)		
		if page.RemainingUsage < 1 || age > CacheMaxAge {
			//fmt.Printf("handler path=%s cache expired\n", path)
			log.Printf("handler path=%s cache expired\n", path)			
			delete(pageCache, path)
			page = nil
		}
	} else {
		//fmt.Printf("handler path=%s cache MISS\n", path)
		log.Printf("handler path=%s cache MISS\n", path)		
	}

	if page == nil {
		//fmt.Printf("handler path=%s refreshing\n", path)
		log.Printf("handler path=%s refreshing\n", path)		
		page, err = loadPage(path)
		if err == nil {
			//fmt.Printf("handler path=%s loaded\n", path)
			log.Printf("handler path=%s loaded\n", path)			
			pageCache[path] = page
		} else {
			//fmt.Printf("handler path=%s not found\n", path)
			log.Printf("handler path=%s not found\n", path)			
			page = pageNotFound(path)
		}
	}

	//fmt.Printf("handler path=%s countdown=%d age=%d/%d\n", path, page.RemainingUsage, cacheAge(page), CacheMaxAge)
	log.Printf("handler path=%s countdown=%d age=%d/%d\n", path, page.RemainingUsage, cacheAge(page), CacheMaxAge)
	
	fmt.Fprintf(w, string(page.Body))
}

func main() {
	pageCache = make(map[string]*Page)

	go cacheTidyPeriodic()

	//fmt.Printf("server starting\n")
	log.Printf("server starting\n")	

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
