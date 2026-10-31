package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func main() {
	url := "https://en.wikipedia.org/wiki/Bear"

	client := &http.Client{Timeout: 15 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("bad request", err)
		return
	}
	req.Header.Set("User-Agent", "bear-crawler/0.1 (+https://github.com/1xxalexx1/bear-crawler)")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("fetch failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("bad status:", resp.Status)
		return
	}

	tokenizer := html.NewTokenizer(resp.Body)
	var links []string
	var wholeText strings.Builder
	// sentences containing the word bear
	var bearSentences []string
	// how many times the word bear appears
	var howManyBears int
	skip := false

	re := regexp.MustCompile(`(?i)\bbears?\b`)
	sentenceSplit := regexp.MustCompile(`[.!?]+\s+`)
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			sentences := sentenceSplit.Split(wholeText.String(), -1)
			for _, s := range sentences {
				if re.MatchString(s) {
					bearSentences = append(bearSentences, s)
				}
			}
			howManyBears = len(re.FindAllString(wholeText.String(), -1))
			fmt.Println(howManyBears)
			fmt.Println(len(bearSentences))
			fmt.Println(len(sentences))
			return
		case html.StartTagToken:
			t := tokenizer.Token()
			if t.Data == "style" || t.Data == "script" {
				skip = true
			}
			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" {
						links = append(links, a.Val)
					}
				}
			}
		case html.TextToken:
			t := tokenizer.Token()
			if skip == false {
				wholeText.WriteString(t.Data)
			}
		case html.EndTagToken:
			t := tokenizer.Token()
			if t.Data == "style" || t.Data == "script" {
				skip = false
			}
		}
	}
}
