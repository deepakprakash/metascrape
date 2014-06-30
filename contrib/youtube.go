package contrib

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/deepakprakash/metascrape/utils"
)

/*
YoutubeHandler implements a handler for Youtube Videos.

A YouTube API Key is required and should be specified as the environment variable
`YOUTUBE_API_KEY`. In the absense of this, matching is not attempted and a `nil, false`
response is returned.

Matching is done if:
  - URL's Host is `www.youtube.com` AND
  - A url query parameter `v` is present which denotes the video ID.
  - A call to the Youtube API endpoint with our data is successful.

Custom return data:
  type: "Video"
  provider: "YouTube"

  extraData:
    datePublished: "Date/Time when this was published."
    duration: "The duration in ISO 8601 date format. (http://en.wikipedia.org/wiki/ISO_8601#Durations)"

    statistics:
      "commentCount"
      "likeCount"
      "dislikeCount"
      "viewCount"
      "favouriteCount"

    creator: "TODO"
    embedDetails: "TODO"
*/
func YoutubeVideoHandler(response *http.Response, doc *goquery.Document) (map[string]interface{}, bool) {

	apiKey := os.Getenv("YOUTUBE_API_KEY")
	// apiKey := "AIzaSyBXcrw5QMbYYrrJoeQGbs7uhGXJoHq6MSk"
	if len(apiKey) == 0 {
		// API Key for Youtube is required
		return nil, false
	}

	if canonicalURL, err := url.Parse(utils.ExtractCanonicalURL(doc, response)); err == nil &&
		canonicalURL.Host == "www.youtube.com" {
		// URL is of www.youtube.com

		if videoID := canonicalURL.Query().Get("v"); len(videoID) > 0 {
			// Video ID present

			// YouTube API end point for videos
			apiURL, _ := url.Parse("https://www.googleapis.com:443/youtube/v3/videos")

			// Add our custom parameters
			params := apiURL.Query()
			params.Add("key", apiKey)
			params.Add("part", "snippet,contentDetails,statistics,player")
			params.Add("fields", "items(id,snippet/title,snippet/publishedAt,snippet/thumbnails/medium,contentDetails/duration,statistics,player/embedHtml)")
			params.Add("id", videoID)

			// Reinit the apiURL
			apiURL.RawQuery = params.Encode()

			if apiResponse, err := http.Get(apiURL.String()); err == nil && apiResponse.StatusCode == http.StatusOK {
				if body, err := ioutil.ReadAll(apiResponse.Body); err == nil {

					type Result struct {
						Items []struct {
							Id             string `json:"id"`
							ContentDetails struct {
								Duration string `json:"duration"`
							} `json:"contentDetails"`
							Snippet struct {
								Title       string    `json:"title"`
								PublishedAt time.Time `json:"publishedAt"`
								Thumbnails  struct {
									Medium struct {
										Height int    `json:"height"`
										Width  int    `json:"width"`
										Url    string `json:"url"`
									} `json:"medium"`
								} `json:"thumbnails"`
							} `json:"snippet"`
							Statistics struct {
								CommentCount   string `json:"commentCount"`
								DislikeCount   string `json:"dislikeCount"`
								FavouriteCount string `json:"favouriteCount"`
								LikeCount      string `json:"likeCount"`
								ViewCount      string `json:"viewCount"`
							} `json:"statistics"`
						} `json:"items"`
					}

					apiData := new(Result)

					if err := json.Unmarshal(body, apiData); err == nil {
						meta, _ := GenericHandler(response, doc)

						if len(apiData.Items) > 0 {
							item := apiData.Items[0]

							extraData := make(map[string]interface{})
							extraData["duration"] = item.ContentDetails.Duration
							extraData["statistics"] = item.Statistics
							extraData["datePublished"] = item.Snippet.PublishedAt

							meta["type"] = "Video"
							meta["provider"] = "YouTube"
							meta["title"] = item.Snippet.Title
							meta["thumbnailUrl"] = item.Snippet.Thumbnails.Medium.Url
							meta["extraData"] = extraData

							return meta, true
						}

					}
				}
			}
		}
	}

	return nil, false
}
