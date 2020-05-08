package versionsfile

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/xabierlaiseca/gowrap/pkg/util/customerrors"
)

type platformGoArchive struct {
	GoArchive `json:",inline"`

	OS   string `json:"os,omitempty"`
	ARCH string `json:"arch,omitempty"`
}

type remoteVersionsFile struct {
	versions map[string][]platformGoArchive
}

// GetArchivesFor available archives indexed by version for the provided CPU
// architecture and OS.
func (rvf *remoteVersionsFile) getArchivesFor(arch, os string) map[string]platformGoArchive {
	foundArchives := make(map[string]platformGoArchive)
	for version, archives := range rvf.versions {
		for _, archive := range archives {
			if archive.ARCH == arch && archive.OS == os {
				foundArchives[version] = archive
				break
			}
		}
	}

	return foundArchives
}

const versionsFileURL = "https://raw.githubusercontent.com/xabierlaiseca/gowrap/master/data/versions.json"

func download() (*remoteVersionsFile, error) {
	response, err := http.Get(versionsFileURL)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, customerrors.Errorf("failed downloading versions file, unexpected status: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	versions := make(map[string][]platformGoArchive)
	err = json.Unmarshal(body, &versions)
	if err != nil {
		return nil, err
	}

	rvf := remoteVersionsFile{
		versions: versions,
	}

	return &rvf, nil
}

func Generate(outputPath string) error {
	downloadsPageDoc, err := getDownloadsPage()
	if err != nil {
		return err
	}

	rvf := extractRemoteVersionsFile(downloadsPageDoc)

	versionsBytes, err := json.MarshalIndent(rvf.versions, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outputPath, versionsBytes, 0600)
}

const downloadsPageURL = "https://golang.org/dl/"

func getDownloadsPage() (*goquery.Document, error) {
	response, err := http.Get(downloadsPageURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, customerrors.Errorf("unexpected status code (%d) while getting downloads page", response.StatusCode)
	}

	return goquery.NewDocumentFromReader(response.Body)
}

var validGoVersionRegex = regexp.MustCompile(`^go[0-9]+(\.[0-9]+){1,2}$`)

func extractRemoteVersionsFile(doc *goquery.Document) remoteVersionsFile {
	versions := make(map[string][]platformGoArchive)

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

	return remoteVersionsFile{versions: versions}
}

var archiveFileRegex = regexp.MustCompile(`^go(?:[0-9]+\.){2,3}([^-]+)-([^\.]+)\..*$`)

func extractArchives(versionSelection *goquery.Selection) []platformGoArchive {
	checksumAlgorithmTitle := versionSelection.
		Find("thead tr th:nth-child(6)").
		Text()

	checksumAlgorithm := strings.SplitN(checksumAlgorithmTitle, " ", 2)[0]

	var archives []platformGoArchive
	versionSelection.Find("tbody tr").
		Filter(`:has(td:contains("archive"))`).
		Each(func(_ int, archiveRowSelection *goquery.Selection) {
			archives = append(archives, extractArchive(archiveRowSelection, checksumAlgorithm))
		})

	return archives
}

func extractArchive(archiveRowSelection *goquery.Selection, checksumAlgorithm string) platformGoArchive {
	link, _ := archiveRowSelection.Find(`td:first-child a`).First().Attr("href")
	filename := path.Base(link)
	matches := archiveFileRegex.FindStringSubmatch(filename)

	checksum := archiveRowSelection.Find(`td:nth-child(6) tt`).Text()
	return platformGoArchive{
		GoArchive: GoArchive{
			URL:               link,
			Checksum:          checksum,
			ChecksumAlgorithm: checksumAlgorithm,
		},
		ARCH: matches[2],
		OS:   matches[1],
	}
}
