package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"encoding/csv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func ExampleScrape(url string, titleChan chan []string, wg *sync.WaitGroup) {
	defer wg.Done()

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("#details-page").Each(func(i int, s *goquery.Selection) {
		title :=  strings.TrimSpace(s.Find(".details-book-main-info__header h1").Text())
		category :=  strings.TrimSpace(s.Find(".details-book-info__content-category a").Text())

		fmt.Printf(">>>> category %s\n",category)

		titleChan <- []string{title, url, category}

	})
}

func main() {
	baseURL := "https://www.rokomari.com/book/"
	titleChan := make(chan []string, 100)
	var wg sync.WaitGroup

	fName := "books.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"title", "url", "category"})


	bookURLStartIndex := 70000
	bookURLEndIndex := 70999
	for n := bookURLStartIndex; n <= bookURLEndIndex; n++ {
		fmt.Println("requesting to : " + baseURL + strconv.Itoa(n))
		wg.Add(1)
		go ExampleScrape(baseURL + strconv.Itoa(n), titleChan, &wg)
	}

	go func() {
		for item := range titleChan {
			// fmt.Println(item)
			writer.Write(item)
		}
	}()

	wg.Wait()
}
