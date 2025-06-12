package main

import (
	"fmt"
	"sync"
	"time"
	"github.com/jackc/pgx/v5"
	"context"
	"log"
	"os"
)

func main() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL not found/set")
	}

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("failed to connect to database :%v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	fmt.Println("Postgres Database is ready")

	return

	startNow := time.Now()

	maxConcurrency := 10

	crawler := CrawlerQueue{
		mu:        &sync.Mutex{},
		wg:        &sync.WaitGroup{},
		elements:  []string{},
		PageURLs:  make(map[uint64]bool),
		urlsFound: 0,
		maxPages:  10000,
	}
	fmt.Println("Initialied CrawlerQueue")

	seed := "https://en.wikipedia.org/wiki/Main_Page"
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
