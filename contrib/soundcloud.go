package contrib

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/deepakprakash/metascrape/utils"
)

/*
SoundcloudAudioHandler implements a basic handler for SoundCloud Audios.

A SoundCloud API Key is required and should be specified as the environment variable
`SOUNDCLOUD_API_KEY`. In the absense of this, matching is not attempted and a `nil, false`
response is returned.

Matching is done if:
  - URL's Host is `soundcloud.com` AND
  - Meta tag exits with `property=og:type` and `content=soundcloud:sound`.

Custom return data:
  type: "Audio"
  provider: "SoundCloud"

  extraData:
    datePublished: "Date/Time when this was published."
    genre: "Genre of the audio."
    duration: "The duration in ISO 8601 date format. (http://en.wikipedia.org/wiki/ISO_8601#Durations)"

    statistics:
      "commentCount"
      "viewCount"
      "favouriteCount"

    creator: "TODO"
    embedDetails: "TODO"
*/
func SoundCloudAudioHandler(response *http.Response, doc *goquery.Document) (map[string]interface{}, bool) {

	apiKey := os.Getenv("SOUNDCLOUD_API_KEY")
	if len(apiKey) == 0 {
		// API Key for SoundCloud is required
		return nil, false
	}

	if canonicalURL, err := url.Parse(utils.ExtractCanonicalURL(doc, response)); err == nil &&
		canonicalURL.Host == "soundcloud.com" {

		if ogType, exists := doc.Find("meta[property='og:type']").First().Attr("content"); exists == true &&
			ogType == "soundcloud:sound" {
			// Detected as SoundCloud Audio URL

			// SoundCloud API end point (this resolves automatically to the required endpoint)
			apiURL, _ := url.Parse("http://api.soundcloud.com/resolve.json")

			// Add our custom parameters
			params := apiURL.Query()
			params.Add("client_id", apiKey)
			params.Add("url", canonicalURL.String())

			// Reinit the apiURL
			apiURL.RawQuery = params.Encode()

			apiResponse, err := http.Get(apiURL.String())

			if err == nil && apiResponse.StatusCode == http.StatusOK {
				if body, err := ioutil.ReadAll(apiResponse.Body); err == nil {

					type Result struct {
						Id          int64  `json:"id"`
						Genre       string `json:"genre"`
						Duration    int64  `json:"duration"`
						Title       string `json:"title"`
						Description string `json:"description"`
						CreatedAt   string `json:"created_at"`
						ArtworkUrl  string `json:"artwork_url"`

						CommentCount   int64 `json:"comment_count"`
						FavouriteCount int64 `json:"favoritings_count"`
						ViewCount      int64 `json:"playback_count"`
					}

					type Statistics struct {
						CommentCount   int64 `json:"commentCount"`
						FavouriteCount int64 `json:"favouriteCount"`
						ViewCount      int64 `json:"viewCount"`
					}

					apiData := new(Result)

					if err := json.Unmarshal(body, apiData); err == nil {
						meta, _ := GenericHandler(response, doc)

						extraData := make(map[string]interface{})

						// Statistics
						stats := Statistics{apiData.CommentCount, apiData.FavouriteCount, apiData.ViewCount}
						extraData["statistics"] = stats

						// Date is not in standard format - so process it.
						dateString := strings.Replace(apiData.CreatedAt, "/", "-", -1)
						dateString = strings.Replace(dateString, " ", "T", 1)
						dateString = strings.Replace(dateString, " +00", "+00:", 1)
						extraData["datePublished"], _ = time.Parse(time.RFC3339Nano, dateString)

						// TODO: Convert duration in seconds to ISO_8601 format.
						extraData["duration"] = apiData.Duration

						extraData["genre"] = apiData.Genre

						meta["type"] = "Audio"
						meta["provider"] = "SoundCloud"

						meta["title"] = apiData.Title
						meta["description"] = apiData.Description
						meta["thumbnailUrl"] = apiData.ArtworkUrl
						meta["extraData"] = extraData

						return meta, true

					}
				}
			}
		}
	}

	return nil, false
}
