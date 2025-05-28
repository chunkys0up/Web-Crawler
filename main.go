package main

import (
	"fmt"
	"sync"
	"net/http"
	"io"
)

/*
	
	"strings"
	"time"
	"context"
	"log"
	"golang.org/x/net/html"
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
		return "Error: " + err.Error()
	}
	defer resp.Body.Close()

	// turns to bytes that need to be turned into a string
	body, err := io.ReadAll(resp.Body)

	
	return string(body)
}

// parse webpage content
func (q *Queue) parseWebPage(body string) {
	
}

func main() {
	queue := Queue{}
	seed := "https://www.wikipedia.org/"

	queue.enQueue((seed))
	
	fmt.Println(queue.fetchPage(queue.deQueue()))

}
