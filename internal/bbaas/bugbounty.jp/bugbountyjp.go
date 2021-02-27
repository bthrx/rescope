//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2021 root4loot
//

package bugbountyjp

import (
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"

	errors "github.com/root4loot/rescope/internal/bbaas/pkg/errors"
	req "github.com/root4loot/rescope/internal/bbaas/pkg/request"
)

// Scrape returns a string containing scope that was scraped from the given program on bugbounty.jp
func Scrape(url string) string {
	var scope []string
	re := regexp.MustCompile(`([\w-]+)\/([\w-]+$)`)
	match := re.FindStringSubmatch(url)
	program := match[2]
	endpoint := "https://bugbounty.jp/program/" + program

	// GET request to endpoint
	respBody, _ := req.GET(endpoint)

	// parse response body to xQuery doc
	doc, _ := htmlquery.Parse(strings.NewReader(respBody))

	// xQuery to grab scope section
	resp := htmlquery.Find(doc, "//dt[contains(text(), 'Scope')]/following-sibling::dd[@class='targetDesc']")

	// get scope contents
	if resp != nil {
		for _, item := range resp {
			scope = append(scope, htmlquery.InnerText(item))
		}
	} else {
		errors.NoScope(url)
	}

	return strings.Join(scope, "\n")
}
