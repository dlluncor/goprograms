// package crawl. Shameless copy of https://gist.github.com/technoweenie
package crawl

import (
	"fmt"
	"github.com/moovweb/gokogiri"
	"io/ioutil"
	"net/http"
	"net/url"
)
 
type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}
 
// crawlUrl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func crawlUrl(url string, depth int, fetcher Fetcher, finished chan bool) {
	if depth <= 0 {
		finished <- false
		return
	}
	_, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		finished <- false
		return
	}
 
	urlCount := len(urls)
	if urlCount > 0 {
		fmt.Printf("found: %s %d in depth: %d\n", url, urlCount, depth)
	}
 
	innerFinished := make(chan bool)
	for _, u := range urls {
		go crawlUrl(u, depth-1, fetcher, innerFinished)
	}
 
	for i := 0; i < urlCount; i += 1 {
		<-innerFinished
	}
 
	finished <- true
 
	return
}
 
type urlFetcher map[string]*fetcherResult
 
type fetcherResult struct {
	body string
	urls []string
}
 
func (f *urlFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := (*f)[url]; ok {
		return res.body, res.urls, nil
	}
 
	return "body", f.scanForUrls(url), nil
}
 
func (f *urlFetcher) scanForUrls(url string) []string {
	body, _ := f.download(url)
	urlStrings := f.urlsInBody(body)
 
	return f.validUrls(url, urlStrings)
}
 
func (f *urlFetcher) validUrls(parentUrl string, urls []string) []string {
	parent, _ := url.Parse(parentUrl)
	validUrls := make([]string, len(urls))
	validUrlCount := 0
 
	for _, urlString := range urls {
		uri, err := url.Parse(urlString)
 
		if err == nil && (uri.Scheme == "http" || uri.Scheme == "https") {
			validUrls[validUrlCount] = parent.ResolveReference(uri).String()
			validUrlCount += 1
		}
	}
 
	return validUrls[:validUrlCount]
}
 
func (f *urlFetcher) urlsInBody(body []byte) []string {
	doc, _ := gokogiri.ParseHtml(body)
	defer doc.Free()
 
	nodes, _ := doc.Search("//a/@href")
	urls := make([]string, len(nodes))
 
	for n, node := range nodes {
		urls[n] = node.String()
	}
 
	return urls
}
 
func (f *urlFetcher) download(url string) ([]byte, error) {
	var resp, err = http.Get(url)
	if err != nil {
		return nil, err
	}
 
	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	return body, err2
}
