package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

/*

"sync"
"time"
*/

func main() {

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL not found/set")
	}

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("failed to connect to database : %v\n", err)
	} else {
		fmt.Println("Connected to Postgres database")

	}

	defer conn.Close(context.Background())

	// Create the table if it doesn't exist
	_, err = conn.Exec(context.Background(), `
		DROP TABLE IF EXISTS articles;
		CREATE TABLE articles (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			url TEXT NOT NULL
		);
	`)
	if err != nil {
		log.Fatalf("failed to create table: %v\n", err)
	} else {
		fmt.Println("Articles table created")
	}

	// test inserting 
	insertArticle(conn, "random_name", "www.randomURL.com")
	fmt.Println("Insert successful!")

	// print articles
	printArticles(conn)

}

func insertArticle(conn *pgx.Conn, name, url string) {
	_, err := conn.Exec(context.Background(),
		"INSERT INTO articles(name, url) VALUES($1, $2)", name, url)
	if err != nil {
		log.Fatalf("error inserting article: %v\n", err)
	}
	fmt.Printf("Inserted article: %s - %s\n", name, url)
}

func printArticles(conn *pgx.Conn) {
	fmt.Println("Fetching articles:")
	rows, err := conn.Query(context.Background(), "SELECT id, name, url FROM articles")
	if err != nil {
		log.Fatalf("select error: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, url string
		err := rows.Scan(&id, &name, &url)
		if err != nil {
			log.Printf("row scan error: %v\n", err)
			continue
		}
		fmt.Printf("ID: %d | Name: %s | URL: %s\n", id, name, url)
	}
}

/*
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
*/
