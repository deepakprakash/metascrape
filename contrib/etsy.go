package contrib

import (
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"

	"github.com/deepakprakash/metascrape/lib"
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

  attributes:
    price: "Numeric value of the list price."
    currency: "The currency (in 3-letter ISO 4217 format) of the price"

    availability: "TODO"
    seller: "TODO"
*/
func EtsyProductHandler(response *http.Response, doc *goquery.Document) (*lib.Metadata, bool) {

	if canonicalURL, err := url.Parse(utils.ExtractCanonicalURL(doc, response)); err == nil && canonicalURL.Host == "www.etsy.com" {

		if price, exists := doc.Find("meta[property='etsymarketplace:price_value']").First().Attr("content"); exists == true {

			if priceCurrency, exists := doc.Find("meta[property='etsymarketplace:currency_code']").First().Attr("content"); exists == true {

				meta, _ := GenericHandler(response, doc)

				meta.SetType("Product")
				meta.SetProvider("Etsy")
				meta.SetAttr("price", price)
				meta.SetAttr("priceCurrency", priceCurrency)

				return meta, true
			}
		}
	}
	return nil, false
}
