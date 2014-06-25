package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/pilu/fresh/runner/runnerutils"

	"github.com/deepakprakash/metascrape"
)

var scraper *metascrape.MetaScraper

func init() {
	scraper = metascrape.Default()
}

func main() {
	// Create a new router object
	router := httprouter.New()

	// Set up the routes
	router.GET("/", getIndex)           // Index
	router.GET("/api/meta", getAPIMeta) // API end point for Meta

	// Setup negroni
	n := negroni.Classic()

	// Add the required middleware
	n.Use(negroni.HandlerFunc(runnerMiddleware)) // Middleware for live reload in dev mode

	// Make negroni use the router
	n.UseHandler(router)

	// TOOD: Convert this into a CLI app with support for extracting the various
	// parameters such as port number, api end point, etc from ENV variables or
	// CLI params.
	// Run the server
	n.Run(":4000")
}

/** MIDDLEWARE **/

/*
runnerMiddleware is the middleware to be used for recompiling and reloading the
server using fresh [http://github.com/pilu/fresh] whenever source changes.
This should be used only in development mode.
*/
func runnerMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if runnerutils.HasErrors() {
		runnerutils.RenderError(w)
	}

	next(w, r)
}

/** HTTP HANDLERS **/

/*
GetIndex handles the requests to the root URL. For now simply replies asking for
a request to the API end point.
*/
func getIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Please use the API end point.")
}

/*
getAPIMeta handles the Meta API requests.

The URL to be queried should be sent as a URL parameter named `url`.

The response `Content-Type` header is set to `application/json`.

On successfully extracting meta data for the given URL, response is sent as json.

Error responses are returned for various situations, including those related to
querying and parsing the given URL.

*/
func getAPIMeta(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	if urlInput := r.FormValue("url"); len(urlInput) > 0 {
		// Call the metascrape lib and write output or error
		if data, err := scraper.Scrape(urlInput); err == nil {
			// Convert to JSON
			jsonResp, _ := json.Marshal(data)

			// Write to HTML
			fmt.Fprint(w, string(jsonResp))

		} else {
			//
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(w, "`url` parameter is empty or missing.", http.StatusBadRequest)
	}
}
