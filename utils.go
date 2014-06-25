package metascrape

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func fetchURL(url string) (*http.Response, error) {

	// TODO: Implement timeouts / data restrictions
	return http.Get(url)
}

func extractTitle(doc *goquery.Document) string {

	if title := doc.Find("title").First().Text(); len(title) > 0 {
		return title
	}

	if title, exists := doc.Find("meta[name='title']").First().Attr("content"); exists == true {
		return title
	}

	if title, exists := doc.Find("meta[property='og:title']").First().Attr("content"); exists == true {
		return title
	}

	if title, exists := doc.Find("meta[name='twitter:title']").First().Attr("content"); exists == true {
		return title
	}

	return ""
}

func extractDescription(doc *goquery.Document) string {

	if description, exists := doc.Find("meta[name='description']").First().Attr("content"); exists == true {
		return description
	}

	if description, exists := doc.Find("meta[property='og:description']").First().Attr("content"); exists == true {
		return description
	}

	if description, exists := doc.Find("meta[name='twitter:description']").First().Attr("content"); exists == true {
		return description
	}

	return ""
}

func extractThumbnailURL(doc *goquery.Document) string {

	if thumbnail, exists := doc.Find("meta[name='thumbnail']").First().Attr("content"); exists == true {
		return thumbnail
	}

	if thumbnail, exists := doc.Find("meta[property='og:image']").First().Attr("content"); exists == true {
		return thumbnail
	}

	if thumbnail, exists := doc.Find("meta[name='twitter:image']").First().Attr("content"); exists == true {
		return thumbnail
	}

	return ""
}

func extractCanonicalURL(doc *goquery.Document, response *http.Response) string {

	if canonical, exists := doc.Find("link[rel='canonical']").First().Attr("href"); exists == true {
		return canonical
	}

	if canonical, exists := doc.Find("meta[property='og:url']").First().Attr("content"); exists == true {
		return canonical
	}

	if canonical, exists := doc.Find("meta[name='twitter:url']").First().Attr("content"); exists == true {
		return canonical
	}

	// No canonical URL found - so return the URL (the latest redirected URL)
	// response.Request.URL contains the final URL that was fetched after redirects
	return response.Request.URL.String()
}
