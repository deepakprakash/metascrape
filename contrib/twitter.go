package contrib

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/PuerkitoBio/goquery"

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
func TwitterProfileHandler(response *http.Response, doc *goquery.Document) (map[string]interface{}, bool) {

	if api != nil {
		if canonicalURL, err := url.Parse(utils.ExtractCanonicalURL(doc, response)); err == nil &&
			canonicalURL.Host == "twitter.com" {

			paths := strings.Split(canonicalURL.Path, "/")

			if len(paths) == 2 {
				username := paths[1]

				if user, err := api.GetUsersShow(username, nil); err == nil {
					meta, _ := GenericHandler(response, doc)

					extraData := make(map[string]interface{})
					meta["extraData"] = extraData

					// Populate the data from user object
					extraData["handle"] = user.ScreenName
					extraData["name"] = user.Name
					extraData["location"] = user.Location
					extraData["bio"] = user.Description
					extraData["dateCreated"], _ = time.Parse(time.RubyDate, user.CreatedAt)

					// Populate the statistics
					stats := make(map[string]interface{})
					extraData["statistics"] = stats

					stats["followingCount"] = user.FriendsCount
					stats["followerCount"] = user.FollowersCount
					stats["tweetCount"] = user.StatusesCount
					stats["favouriteCount"] = user.FavouritesCount

					// Assign the general properties
					meta["thumbnailUrl"] = user.ProfileImageUrlHttps
					meta["type"] = "Profile"
					meta["provider"] = "Twitter"

					return meta, true
				}
			}
		}
	}

	return nil, false
}
