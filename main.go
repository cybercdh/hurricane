package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// fetchURL fetches the URL and returns the root HTML node.
func fetchURL(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// parseHTML recursively searches for <a> tags with hrefs that start with "/net/".
func parseHTML(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			if a.Key == "href" && strings.HasPrefix(a.Val, "/net/") {
				ipRange := strings.TrimPrefix(a.Val, "/net/")
				fmt.Println(ipRange)
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		parseHTML(c)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: hurricane <search query>")
		return
	}

	searchQuery := os.Args[1]
	// URL encode the search query to ensure it is safe to include in a URL
	encodedQuery := url.QueryEscape(searchQuery)

	// Insert the encoded search query into the URL
	url := fmt.Sprintf("https://bgp.he.net/search?search%%5Bsearch%%5D=%s&commit=Search", encodedQuery)

	doc, err := fetchURL(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}

	parseHTML(doc)
}
