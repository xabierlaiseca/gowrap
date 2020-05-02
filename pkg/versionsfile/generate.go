package versionsfile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GoArchive struct {
	OS     string
	CPU    string
	URL    string
	SHA256 string
}

func Generate(outputPath string) error {
	downloadsPageDoc, err := getDownloadsPage()
	if err != nil {
		return err
	}

	versions := extractVersions(downloadsPageDoc)

	versionsBytes, err := json.Marshal(versions)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outputPath, versionsBytes, 0644)
}

const downloadsPageURL = "https://golang.org/dl/"

func getDownloadsPage() (*goquery.Document, error) {
	response, err := http.Get(downloadsPageURL)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code (%d) while getting downloads page", response.StatusCode)
	}

	return goquery.NewDocumentFromReader(response.Body)
}

var validGoVersionRegex = regexp.MustCompile(`^go[0-9]+(\.[0-9]+){1,2}$`)

func extractVersions(doc *goquery.Document) map[string][]GoArchive {
	versions := make(map[string][]GoArchive)

	doc.Find(`div[id^="go"]`).
		FilterFunction(func(_ int, selection *goquery.Selection) bool {
			id, _ := selection.Attr("id")
			return validGoVersionRegex.MatchString(id)
		}).
		Each(func(_ int, selection *goquery.Selection) {
			id, _ := selection.Attr("id")
			version := strings.TrimPrefix(id, "go")
			archives := extractArchives(selection.Find("table"))
			versions[version] = archives
		})

	return versions
}

var archiveFileRegex = regexp.MustCompile(`^go(?:[0-9]+\.){2,3}([^-]+)-([^\.]+)\..*$`)

func extractArchives(versionSelection *goquery.Selection) []GoArchive {
	var archives []GoArchive

	versionSelection.Find("tbody tr").
		Filter(`:has(td:contains("archive"))`).
		Each(func(_ int, archiveRowSelection *goquery.Selection) {
			archives = append(archives, extractArchive(archiveRowSelection))
		})

	return archives
}

func extractArchive(archiveRowSelection *goquery.Selection) GoArchive {
	link, _ := archiveRowSelection.Find(`td:first-child a`).First().Attr("href")
	filename := path.Base(link)
	matches := archiveFileRegex.FindStringSubmatch(filename)

	checksum := archiveRowSelection.Find(`td:nth-child(6) tt`).Text()
	return GoArchive{
		CPU:    matches[2],
		OS:     matches[1],
		URL:    link,
		SHA256: checksum,
	}
}
