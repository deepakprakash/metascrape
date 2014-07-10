package metascrape

import (
	"github.com/deepakprakash/metascrape/contrib"
	"github.com/deepakprakash/metascrape/lib"
)

func New() *lib.MetaScraper {
	scraper := new(lib.MetaScraper)
	// scraper.handlers = []lib.ScrapeHandler{}

	return scraper
}

func Default() *lib.MetaScraper {
	scraper := New()

	scraper.Use(contrib.GenericHandler)
	scraper.Use(contrib.EtsyProductHandler)
	scraper.Use(contrib.YouTubeVideoHandler)
	scraper.Use(contrib.SoundCloudAudioHandler)
	scraper.Use(contrib.TwitterProfileHandler)
	scraper.Use(contrib.TwitterStatusHandler)

	return scraper
}
