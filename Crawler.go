package main

import (
	"golang.org/x/net/html"
	"github.com/jackc/pgx/v5"
	"hash/fnv"
	"io"
	"net/http"
	"strings"
	"sync"
)

type CrawlerQueue struct {
	mu        *sync.Mutex
	wg        *sync.WaitGroup
	elements  []string
	PageURLs  map[uint64]bool // true / false if a hashed url is found
	urlsFound uint64
	maxPages  uint64
	deQueued  uint64
}

func (crawler *CrawlerQueue) size() uint64 {
	crawler.mu.Lock()
	defer crawler.mu.Unlock()

	return crawler.urlsFound
}

func hash(url string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(url))

	return hash.Sum64()
}

func (crawler *CrawlerQueue) addToSet(url string) bool {
	crawler.mu.Lock()
	defer crawler.mu.Unlock()

	hashed := hash(url)
	if crawler.PageURLs[hashed] {
		return false
	}

	crawler.PageURLs[hashed] = true
	crawler.urlsFound++
	return true
}

func (crawler *CrawlerQueue) enQueue(url string) {
	crawler.mu.Lock()
	defer crawler.mu.Unlock()

	crawler.elements = append(crawler.elements, url)
}

func (crawler *CrawlerQueue) deQueue() (bool, string) {
	crawler.mu.Lock()
	defer crawler.mu.Unlock()

	if len(crawler.elements) > 0 {
		val := crawler.elements[0]
		crawler.elements = crawler.elements[1:]
		crawler.deQueued++
		return true, val
	}

	return false, "failed"
}

func getURLs(url string) ([]string, bool) {
	urls := []string{}

	// get the html body
	resp, err := http.Get(url)
	if err != nil {
		return urls, false
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return urls, false
	}

	// we have to parse through the html body and get urls

	htmlBody := string(body)
	tokenIndex := html.NewTokenizer(strings.NewReader(htmlBody))

	for {
		// exit and return the urls we find
		if tokenIndex.Next() == html.ErrorToken {
			return urls, true
		}

		token := tokenIndex.Token()
		if token.Type == html.StartTagToken && token.Data == "a" {
			// now we have to check if the html line has a valid url

			for _, attr := range token.Attr {
				// make sure its a wikipedia article
				if attr.Key == "href" {
					href := attr.Val
					// what href shows: /wiki/Mongoose#Taxonomy

					if strings.HasPrefix(href, "/wiki/") && !strings.Contains(href, ":") {
						fullURL := "https://en.wikipedia.org" + href
						urls = append(urls, fullURL)
					}
				}
			}
		}
	}
}

func (crawler *CrawlerQueue) crawl(conn *pgx.Conn) {
	defer crawler.wg.Done()

	for {
		// maximum pages found
		if crawler.size() > crawler.maxPages {
			return
		}

		err, url := crawler.deQueue()

		if !err {
			//log.Fatal("Failed to deQueue from queue\n")
			return
		}

		// Fetching HTML, and content
		urls, ok := getURLs(url)
		if !ok {
			//fmt.Printf("Failed to get urls. Skipping...\n")
			continue
		}

		// get rid of duplicate urls and then enQueue them
		for _, newURL := range urls {
			
			// skips found urls
			if !crawler.addToSet(newURL) {
				continue
			}
			insertArticle(conn, newURL[len("https://en.wikipedia.org/wiki/"):], newURL)
			crawler.enQueue(newURL)
		}
	}
}

