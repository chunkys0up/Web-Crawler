package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	startNow := time.Now()

	maxConcurrency := 10

	crawler := CrawlerQueue{
		mu:        &sync.Mutex{},
		wg:        &sync.WaitGroup{},
		elements:  []string{},
		PageURLs:  make(map[uint64]bool),
		urlsFound: 0,
		maxPages:  25000,
	}
	fmt.Println("Initialied CrawlerQueue")

	seed := "https://en.wikipedia.org/wiki/Dog"
	crawler.enQueue(seed)
	crawler.addToSet(seed)

	// infinite loop to crawl the web in batches
	fmt.Printf("Creating goroutines\n")
	for {
		for i := 0; i < maxConcurrency; i++ {
			
			crawler.wg.Add(1)
			go crawler.crawl()
		}

		crawler.wg.Wait()

		if crawler.size() >= crawler.maxPages {
			break
		}
	}

	fmt.Println("Completed web crawl")
	fmt.Println("----- Crawl Stats -----")
	fmt.Println("URLs Found:", crawler.size())
	fmt.Println("URLs deQueued:", crawler.deQueued)
	fmt.Println("Time: ", time.Since(startNow))
}
