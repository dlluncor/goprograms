package crawl

import(
  "fmt"
  "time"
)

func start() { 
	fetcher := new(urlFetcher)
 
	t0 := time.Now()
	finished := make(chan bool)
	go crawlUrl("http://golang.org/", 3, fetcher, finished)
	<-finished
	t1 := time.Now()
	fmt.Printf("%v\n", t1.Sub(t0))
}

func Crawl() {
  fmt.Printf("About to crawl apps. \n")
  start()
}
