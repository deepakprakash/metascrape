package contrib

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"

	"github.com/deepakprakash/metascrape/lib"
	"github.com/deepakprakash/metascrape/utils"
)

func GenericHandler(response *http.Response, doc *goquery.Document) (*lib.Metadata, bool) {
	meta := lib.NewMetadata()

	meta.SetType("Webpage")
	meta.SetProvider("")

	meta.SetAttr("title", utils.ExtractTitle(doc))
	meta.SetAttr("description", utils.ExtractDescription(doc))
	meta.SetAttr("thumbnailUrl", utils.ExtractThumbnailURL(doc))
	meta.SetAttr("url", utils.ExtractCanonicalURL(doc, response))

	return meta, true
}
