package main

import (
	"fmt"
	"sync"
	"net/http"
	"io"
	"strings"
	"golang.org/x/net/html"
	"log"
)	

/*
	
	"time"
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

func (q *Queue) fetchPage(url string) string {
	q.mu.Lock()
	defer q.mu.Unlock()

	resp, err := http.Get(url)

	//meaning that there's an error and its not blank
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// turns to bytes that need to be turned into a string
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

// parse webpage content
func (q *Queue) parseWebPage(htmlBody string) {
	
	reader := strings.NewReader(htmlBody)
	doc, err := html.Parse(reader)

	if err != nil {
		log.Fatal(err)
	}

	// loop through eac html node
	for n := range doc.Descendants() {
		// loop through its attributes, html has a key and value like <href = "www.example.com">
		// key = href, value = www.example.com
		for _, a := range n.Attr {
			if n.Type == html.ElementNode && a.Key == "href" {
				fmt.Println(a.Val)
			}
		}	
	}
}

func main() {
	queue := Queue{}
	seed := "https://www.wikipedia.org/"

	queue.enQueue((seed))
	
	htmlBody := queue.fetchPage(queue.deQueue())

	queue.parseWebPage(htmlBody)

}
