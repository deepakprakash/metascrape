Metascrape
==========

Metascrape is a simple, extendable Golang library that helps you extract meaningful metadata from URLs in a structured manner.

A minimal, fast, production ready server that has a JSON API end point is bundled that you can deploy to use Metascrape with its default integrations.

## Why & How

When you are programmatically trying to make sense of URLs or websites, it usually involves two parts

1. Scraping the URLs.

  Relatively simple and has a large variety of libraries and tools to help with it. However, it can end up being a lot of grunt work if simple scraping of the site doesn't give you enough information and you need to resort to site specific APIs.

2. Categorizing and extracting sensible data.
  
  This is much harder, mainly because:
  - There are very few standardized meta tags/representation and there is no guarantee a given website will implement any of it.
  - Data is very different for varied types of URLs - for instance a YouTube Video vs an Amazon Product.
  - Variation in representation of similar data across websites/providers. Eg: Dates, Time, Units of money, Geolocations, etc.

Metascrape (hopefully) abstracts away the complexity of this process by:
- Providing a single end point to query for URL/Website metadata (either using the library or bundled webserver).
- Standardizing as much of the data as possible. For instance common attributes like "datePublished" are often represented in different formats by various websites/providers while Metascrape will always provide it in [ISO 8601](http://en.wikipedia.org/wiki/ISO_8601)format.
- Doing out of the box, specific handling for a large number of common entities on the web. Eg: YouTube Videos, Twitter Profiles, Tweets, SoundCloud Audios, etc. (We start off with a small set and expand as we go.)
- Enabling extention of the library with your own specific integration if you need to handle something thats not built-in or need to handle it differently.

## Uses
- Displaying "rich snippets" for websites like those found on Facebook, Google search results, Nuzzel, etc.
- Extracting useful information from websites for various kinds of analytics, trend analysis, better search, etc.

Note: While explaining Metascrape, some people tend to believe that Metascrape will magically take care of their core analysis, data processing etc. It doesn't. Metascrape simply makes accessing semantic web data easier and structured. You still need to do whatever you need to do with that data.

## Library
If you use Golang for your backend work, the Metascrape library can simply be dropped in and only a function call is necessary to extract data from a URL.

Installation, Usage, Extending the library and other info can be found [here](USAGE.md).

## Built-in Server
If you would rather run Metascrape separately as a web server and have it return the extracted data as JSON via a simple API, please take a look at the server [readme](metascrape/README.md)

## Contributing
Please open a Github issue for any bugs you find. If you have an attached pull request, even better!

- Search the code for "TODO" to check obvious things that need to be done soon.
- We need a LOT more integrations to be added to "contrib" to make Metascrape really useful out of the box. If you can implement integrations for any popular web entity, that would be super great.
- Tests (ahem..)
