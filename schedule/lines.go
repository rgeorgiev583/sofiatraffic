package schedule

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/rgeorgiev583/sofiatraffic/schedule/l10n"

	"golang.org/x/net/html/atom"

	"golang.org/x/net/html"
)

// Lines represents the numbers of all urban transit lines.
type Lines struct {
	BusLineNumbers, TrolleybusLineNumbers, TramLineNumbers []string
}

const (
	linesPagePath = "/"
)

// GetLines fetches and returns all urban transit lines.
func GetLines() (lines *Lines, err error) {
	linesPageURL := &url.URL{
		Scheme: Scheme,
		Host:   Hostname,
		Path:   linesPagePath,
	}
	response, err := http.Get(linesPageURL.String())
	if err != nil {
		err = fmt.Errorf("could not initiate HTTP GET request to the schedule line list page: %s", err.Error())
		return
	}
	defer response.Body.Close()

	lines = &Lines{
		BusLineNumbers:        []string{},
		TrolleybusLineNumbers: []string{},
		TramLineNumbers:       []string{},
	}
	for tokenizer := html.NewTokenizer(response.Body); tokenizer.Next() != html.ErrorToken; {
		token := tokenizer.Token()
		if token.Type == html.StartTagToken && token.DataAtom == atom.A {
			for _, attr := range token.Attr {
				if atom.Lookup([]byte(attr.Key)) == atom.Href {
					lineNumber, err := url.PathUnescape(path.Base(attr.Val))
					if err != nil {
						return lines, err
					}

					if strings.HasPrefix(attr.Val, VehicleTypeTram) {
						lines.TramLineNumbers = append(lines.TramLineNumbers, lineNumber)
					} else if strings.HasPrefix(attr.Val, VehicleTypeTrolleybus) {
						lines.TrolleybusLineNumbers = append(lines.TrolleybusLineNumbers, lineNumber)
					} else if strings.HasPrefix(attr.Val, VehicleTypeBus) {
						lines.BusLineNumbers = append(lines.BusLineNumbers, lineNumber)
					}
				}
			}
		}
	}
	return
}

func (ls *Lines) String() (str string) {
	str += "* " + l10n.Translator[l10n.BusLines] + ": " + strings.Join(ls.BusLineNumbers, ", ") + "\n"
	str += "* " + l10n.Translator[l10n.TrolleybusLines] + ": " + strings.Join(ls.TrolleybusLineNumbers, ", ") + "\n"
	str += "* " + l10n.Translator[l10n.TramLines] + ": " + strings.Join(ls.TramLineNumbers, ", ") + "\n"
	return
}
