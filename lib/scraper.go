package lib

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type ScrapeHandler func(resp *http.Response, doc *goquery.Document) (*Metadata, bool)

type MetaScraper struct {
	handlers []ScrapeHandler
}

func (scraper *MetaScraper) Use(handler ScrapeHandler) {
	scraper.handlers = append(scraper.handlers[:0], append([]ScrapeHandler{handler}, scraper.handlers[0:]...)...)
}

func (scraper *MetaScraper) Scrape(urlInput string) (*Metadata, error) {

	if pURL, err := url.ParseRequestURI(urlInput); err == nil {
		// This is a valid URL.

		if response, err := fetchURL(pURL.String()); err == nil {
			// Able to query URL and get data properly
			// head, _ := getHead(response)

			// Parse the response body and create tree structure required for goquery
			if doc, err := goquery.NewDocumentFromResponse(response); err == nil {

				for _, handler := range scraper.handlers {
					// metaData := GenericHandler(response, doc)
					if metaData, matched := handler(response, doc); matched == true {
						// Handler was able to process
						return metaData, nil
					}
				}

				// None of the handlers were able to match
				return nil, errors.New("None of the handlers were able to match and scrape the URL.")

			} else {
				// http.Error(w, "Unable to parse the response.", http.StatusBadRequest)
				return nil, errors.New("Unable to parse the response.")
			}

		} else {
			// http.Error(w, "Unable to fetch URL or similar issue. Please check URL.", http.StatusBadRequest)
			return nil, errors.New("Unable to fetch URL or similar issue. Please check URL.")
		}

	} else {
		// http.Error(w, fmt.Sprint("Invalid `url` input: ", urlInput), http.StatusBadRequest)
		return nil, errors.New(fmt.Sprint("Invalid `url` input: ", urlInput))
	}
}

func fetchURL(url string) (*http.Response, error) {

	// TODO: Implement timeouts / data restrictions
	return http.Get(url)
}
