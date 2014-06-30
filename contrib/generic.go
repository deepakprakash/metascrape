package contrib

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"

	"github.com/deepakprakash/metascrape/utils"
)

func GenericHandler(response *http.Response, doc *goquery.Document) (map[string]interface{}, bool) {
	meta := make(map[string]interface{})

	meta["type"] = "Webpage"
	meta["provider"] = ""

	meta["title"] = utils.ExtractTitle(doc)
	meta["description"] = utils.ExtractDescription(doc)
	meta["thumbnailUrl"] = utils.ExtractThumbnailURL(doc)
	meta["url"] = utils.ExtractCanonicalURL(doc, response)

	meta["extraData"] = nil

	return meta, true
}
