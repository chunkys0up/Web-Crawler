# Web-Crawler

This projects crawls Wikipedia articles and stores their metadata (name and URL) in a PostgreSQL database. Built with Go and Docker Compose.

## Prerequisites
- Docker
- Docker Compose

## To Run the Repository
```bash
git clone https://github.com/chunkys0up/Web-Crawler.git
```

Make sure dependencies are installed
```bash
`go mod tidy`
```

## Run the Docker Compose
```bash
docker-compose up --build
```

Optional Reset the Database
```bash
docker-compose down -v
```

Might have to run `docker-compose up --build` twice in case the database isn't set up before the program tries running.


