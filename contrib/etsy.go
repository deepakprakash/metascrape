package contrib

import (
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"

	"github.com/deepakprakash/metascrape/utils"
)

/*
EtsyProductHandler implements a basic handler for Etsy Products/Listings.

Matching is done if:
  - URL's Host is `www.etsy.com` AND
  - A valid `price` meta data is extracted.

Custom return data:
  type: "Product"
  provider: "Etsy"

  extraData:
    price: "Numeric value of the list price."
    currency: "The currency (in 3-letter ISO 4217 format) of the price"

    availability: "TODO"
    seller: "TODO"
*/
func EtsyProductHandler(response *http.Response, doc *goquery.Document) (map[string]interface{}, bool) {

	if canonicalURL, err := url.Parse(utils.ExtractCanonicalURL(doc, response)); err == nil && canonicalURL.Host == "www.etsy.com" {

		extraData := make(map[string]interface{})

		if price, exists := doc.Find("meta[property='etsymarketplace:price_value']").First().Attr("content"); exists == true {
			extraData["price"] = price
		}

		if priceCurrency, exists := doc.Find("meta[property='etsymarketplace:currency_code']").First().Attr("content"); exists == true {
			extraData["priceCurrency"] = priceCurrency
		}

		if _, exists := extraData["price"]; exists {
			meta, _ := GenericHandler(response, doc)

			meta["type"] = "Product"
			meta["provider"] = "Etsy"
			meta["extraData"] = extraData

			return meta, true
		} else {
			return nil, false
		}

	} else {
		return nil, false
	}
}
