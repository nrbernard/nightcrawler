package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	url := args[0]
	fmt.Printf("starting crawl of: %s\n", url)

	html, err := getHTML(url)
	if err != nil {
		fmt.Printf("error getting HTML: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s", html)
	os.Exit(0)
}
