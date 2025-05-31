package main

import (
	"fmt"
	"golang.org/x/net/html"
	"hash/fnv"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

/*

	"context"

*/

type Queue struct {
	elementsInQueue int
	totalQueued     int
	mu              sync.Mutex
	elements        []string
}

func (q *Queue) enQueue(url string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.elements = append(q.elements, url)
	q.totalQueued++
	q.elementsInQueue++
}

func (q *Queue) deQueue() string {
	q.mu.Lock()
	defer q.mu.Unlock()

	val := "Queue is empty"
	if len(q.elements) > 0 {
		val = q.elements[0]
		q.elements = q.elements[1:]
		q.elementsInQueue--
	}

	return val
}

type HashSet struct {
	length uint64
	set    map[uint64]bool // true/false if a hashed url is found
	mu     sync.Mutex
}

func hash(url string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(url))

	return hash.Sum64()
}

func (h *HashSet) add(url string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	// already exists, then we return false
	hashed := hash(url)
	if h.set[hashed] {
		return false
	}

	h.set[hashed] = true
	h.length++
	return true
}

func (h *HashSet) size() uint64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.length
}

func fetchPage(url string) (string, bool) {
	resp, err := http.Get(url)
	//meaning that there's an error and its not blank
	if err != nil {
		return "", true
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", true
	}

	return string(body), false
}

func getHref(token html.Token) (url string, ok bool) {
	for _, a := range token.Attr {
		if a.Key == "href" {
			if len(a.Val) > 0 && strings.HasPrefix(a.Val, "http") {
				return a.Val, true
			}
			return "", false
		}
	}
	return "", false
}

// parse webpage content
func parseWebPage(htmlBody string, q *Queue, visited *HashSet, maxWebsites int) {
	tokenIndex := html.NewTokenizer(strings.NewReader(htmlBody))
	for {
		if tokenIndex.Next() == html.ErrorToken {
			return
		}

		token := tokenIndex.Token()
		if token.Type == html.StartTagToken && token.Data == "a" {
			url, ok := getHref(token)
			if ok && visited.size() < uint64(maxWebsites) && visited.add(url) {
				q.enQueue(url)
				//fmt.Println(visited.size()," ",url)
			}
		}
	}
}

func main() {
	startNow := time.Now()

	maxWebsites := 1000
	queue := Queue{}
	set := HashSet{set: make(map[uint64]bool)}
	seed :="https://en.wikipedia.org/wiki/Dog"
	queue.enQueue((seed))
	set.add(seed)

	for set.size() < uint64(maxWebsites) {

		 url := queue.deQueue()
		htmlBody, err := fetchPage(url)

		if !err {
			parseWebPage(htmlBody, &queue, &set, maxWebsites)
		}
	}

	fmt.Println("\n----- Web crawler stats -----")
	fmt.Println("set size:", set.size())
	fmt.Println("Current elements in queue:", queue.elementsInQueue)
	fmt.Println("The operation took:",time.Since(startNow))
}

// ----- Web crawler stats -----
// set size: 1000
// Current elements in queue: 994
// The operation took: 2.618358333s
// andrewnguyen@MacBook-Pro-2 WebCrawler % go run main.go

// ----- Web crawler stats -----
// set size: 1000
// Current elements in queue: 994
// The operation took: 2.576759834s
// andrewnguyen@MacBook-Pro-2 WebCrawler % go run main.go

// ----- Web crawler stats -----
// set size: 1000
// Current elements in queue: 994
// The operation took: 2.62936975s