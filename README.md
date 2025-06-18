# Nightcrawler 

A web crawler written in Go. 

## Prerequisites

- Go 1.x or later

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/nightcrawler.git
cd nightcrawler
```

2. Build the project:
```bash
go build
```

## Usage

Run the crawler with a base URL:

```bash
./nightcrawler <base-url> [max-concurrency] [max-pages]
```

### Arguments

- `base-url` (required): The starting URL to crawl
- `max-concurrency` (optional): Maximum number of concurrent crawls (default: 3)
- `max-pages` (optional): Maximum number of pages to crawl (default: 10)