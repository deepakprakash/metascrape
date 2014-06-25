package metascrape
func fetchURL(url string) (*http.Response, error) {

	// TODO: Implement timeouts / data restrictions
	return http.Get(url)
}
