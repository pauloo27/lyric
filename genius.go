// Package lyric can be used to search and fetch song lyrics
package lyric

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/anaskhan96/soup"
)

// Fetch fetchs a lyrics by it's url (the parameter path (not called url to
// avoid problems with the url package)). It returns the lyric and an error.
func Fetch(path string) (lyric string, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", path, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:59.0) Gecko/20100101 Firefox/81.0")

	res, err := client.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode == 404 {
		err = fmt.Errorf("Cannot find lyrics (status code 404)")
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}

	doc := soup.HTMLParse(string(body))

	lyricDiv := doc.Find("div", "class", "lyrics")

	if lyricDiv.Error == nil {
		lyric = strings.TrimSpace(lyricDiv.FullText())
	} else {
		for _, div := range doc.FindAll("div", "class", "eOLwDW") {

			html := strings.ReplaceAll(div.HTML(), "<br/>", "<br/>\n")

			lyric += strings.TrimSpace(soup.HTMLParse(html).FullText()) + "\n"
		}
		lyric = strings.TrimSpace(lyric)
	}

	return
}

// used in the search methods to find a song entry
var lyricURLRe = regexp.MustCompile(`https:\/\/genius.com/[^/]+-lyrics`)

var ErrNotFound = errors.New("Not found")

func Search(query string) (string, error) {
	path := fmt.Sprintf("https://genius.com/api/search/multi?per_page=5&q=%s", url.QueryEscape(query))
	res, err := http.Get(path)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var jsonRes map[string]interface{}

	err = json.Unmarshal(buffer, &jsonRes)
	if err != nil {
		return "", err
	}

	sections := jsonRes["response"].(map[string]interface{})["sections"].([]interface{})

	for _, section := range sections {
		if section.(map[string]interface{})["type"] == "song" {
			for _, hit := range section.(map[string]interface{})["hits"].([]interface{}) {
				return hit.(map[string]interface{})["result"].(map[string]interface{})["url"].(string), nil
			}
		}
	}
	return "", ErrNotFound
}

// SearchDDG searchs for a query using DuckDuckGO. Search engines can deal with
// typos and not exact searchs. DuckDuckGO have a rate limit, so don't call the
// same search too many times. It returns the URL of the lyrics and an error,
// you can fetch the actual lyrics using Fetch(path).
func SearchDDG(query string) (string, error) {
	path := fmt.Sprintf("https://html.duckduckgo.com/html/?q=site:genius.com+%s", url.QueryEscape(query))

	client := &http.Client{}
	req, err := http.NewRequest("GET", path, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:59.0) Gecko/20100101 Firefox/81.0")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	doc := soup.HTMLParse(string(buffer))
	if err != nil {
		return "", err
	}

	for _, result := range doc.FindAll("a", "class", "result__url") {
		r := fmt.Sprintf("https://%s", strings.TrimSpace(result.Text()))

		if lyricURLRe.MatchString(r) {
			return r, nil
		}
	}

	return "", ErrNotFound
}
