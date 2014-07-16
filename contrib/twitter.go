package contrib

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/PuerkitoBio/goquery"

	"github.com/deepakprakash/metascrape/lib"
	"github.com/deepakprakash/metascrape/utils"
)

var (
	apiConsumerKey    = ""
	apiConsumerSecret = ""
	apiAccessToken    = ""
	apiAccessSecret   = ""
)

var api *anaconda.TwitterApi = nil

func init() {
	apiConsumerKey = os.Getenv("TWITTER_API_KEY")
	apiConsumerSecret = os.Getenv("TWITTER_API_SECRET")
	apiAccessToken = os.Getenv("TWITTER_ACCESS_TOKEN")
	apiAccessSecret = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

	if len(apiConsumerKey) > 0 && len(apiConsumerSecret) > 0 &&
		len(apiAccessToken) > 0 && len(apiAccessSecret) > 0 {

		anaconda.SetConsumerKey(apiConsumerKey)
		anaconda.SetConsumerSecret(apiConsumerSecret)
		api = anaconda.NewTwitterApi(apiAccessToken, apiAccessSecret)
	}
}

/*
TwitterProfileHandler implements a handler for Twitter Profiles.

The following environment variables are required to be set since the Twitter API
is internally queried.
`TWITTER_API_KEY`
`TWITTER_API_SECRET`
`TWITTER_ACCESS_TOKEN`
`TWITTER_ACCESS_TOKEN_SECRET`

Matching is done if:
  - URL's Host is `twitter.com` AND
  - URL path matches the general Twitter Profile scheme, ie, twitter.com/<username>
  - The Twitter API call is successful

Custom return data:
  type: "Profile"
  provider: "Twitter"

  extraData:
    name: "Name on Twitter."
    bio: "Twitter bio."
    location: "Text representation of the location."
    handle: "Twitter username."
    dateCreated: "Date/Time when this profile was created."

    statistics:
      "followerCount"
      "followingCount"
      "tweetCount"
      "favouriteCount"

*/
func TwitterProfileHandler(response *http.Response, doc *goquery.Document) (*lib.Metadata, bool) {

	if api != nil {
		if canonicalURL, err := url.Parse(utils.ExtractCanonicalURL(doc, response)); err == nil &&
			canonicalURL.Host == "twitter.com" {

			paths := strings.Split(canonicalURL.Path, "/")

			if len(paths) == 2 {
				username := paths[1]

				if user, err := api.GetUsersShow(username, nil); err == nil {
					meta, _ := GenericHandler(response, doc)

					// Assign the general properties
					meta.SetType("Profile")
					meta.SetProvider("Twitter")
					meta.SetAttr("thumbnailUrl", user.ProfileImageUrlHttps)

					// Populate the data from user object
					meta.SetAttr("handle", user.ScreenName)
					meta.SetAttr("name", user.Name)
					meta.SetAttr("location", user.Location)
					meta.SetAttr("bio", user.Description)
					dateTime, _ := time.Parse(time.RubyDate, user.CreatedAt)
					meta.SetAttr("dateCreated", dateTime)

					// Populate the statistics
					stats := make(map[string]interface{})
					stats["followingCount"] = user.FriendsCount
					stats["followerCount"] = user.FollowersCount
					stats["tweetCount"] = user.StatusesCount
					stats["favoriteCount"] = user.FavouritesCount

					meta.SetAttr("statistics", stats)

					return meta, true
				}
			}
		}
	}

	return nil, false
}

func extractEntities(entities *anaconda.Entities) map[string]interface{} {
	mapEntities := make(map[string]interface{})

	// Handle the hashtags
	hashTags := []string{}
	for _, hashTag := range entities.Hashtags {
		hashTags = append(hashTags, hashTag.Text)
	}
	mapEntities["hashTags"] = hashTags

	// Handle the URLs
	type twitterUrl map[string]interface{}
	urls := []twitterUrl{}

	for _, url := range entities.Urls {
		tUrl := make(twitterUrl)
		tUrl["short"] = url.Url
		tUrl["original"] = url.Expanded_url

		urls = append(urls, tUrl)
	}
	mapEntities["urls"] = urls

	return mapEntities
}

func extractTweetData(tweet *anaconda.Tweet) map[string]interface{} {
	data := make(map[string]interface{})

	// Populate the data from user object
	data["content"] = tweet.Text
	data["datePublished"], _ = time.Parse(time.RubyDate, tweet.CreatedAt)

	// Populate the entities
	data["entities"] = extractEntities(&tweet.Entities)

	// Populate the statistics
	stats := make(map[string]interface{})
	data["statistics"] = stats

	stats["retweetCount"] = tweet.RetweetCount
	stats["favoriteCount"] = tweet.FavoriteCount

	return data
}

/*
TwitterStatusHandler implements a handler for Twitter Statuses.

The following environment variables are required to be set since the Twitter API
is internally queried.
`TWITTER_API_KEY`
`TWITTER_API_SECRET`
`TWITTER_ACCESS_TOKEN`
`TWITTER_ACCESS_TOKEN_SECRET`

Matching is done if:
  - URL's Host is `twitter.com` AND
  - URL path matches the general Twitter Status scheme, ie, twitter.com/<username>/status/<id>
  - The Twitter API call is successful

Custom return data:
  type: "Status"
  provider: "Twitter"

  extraData:
    content: "Full text content of tweet."
    datePublished: "Date/Time when this tweet was created."

    statistics:
      "retweetCount"
      "favouriteCount"
    entities:
      hashTags: [
        // Array of hashtag strings
      ]
      urls: [
        {
          original: // The original URL
          short: // Twitter shortened URL
        },
      ]

    author: "TODO"

*/
func TwitterStatusHandler(response *http.Response, doc *goquery.Document) (*lib.Metadata, bool) {

	if api != nil {
		if canonicalURL, err := url.Parse(utils.ExtractCanonicalURL(doc, response)); err == nil &&
			canonicalURL.Host == "twitter.com" {

			paths := strings.Split(canonicalURL.Path, "/")

			if len(paths) == 4 && paths[2] == "status" {

				if statusId, err := strconv.ParseInt(paths[3], 0, 64); err == nil {

					if tweet, err := api.GetTweet(statusId, nil); err == nil {
						meta, _ := GenericHandler(response, doc)

						data := extractTweetData(&tweet)
						for key, val := range data {
							meta.SetAttr(key, val)
						}

						// Assign the general properties
						// TODO: thumbnailUrl
						meta.SetType("Status")
						meta.SetProvider("Twitter")

						return meta, true
					}
				}

			}
		}
	}

	return nil, false
}
