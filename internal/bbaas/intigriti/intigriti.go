//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2021 root4loot
//

package intigriti

import (
	"regexp"
	"strings"

	errors "github.com/root4loot/rescope/internal/bbaas/pkg/errors"
	req "github.com/root4loot/rescope/internal/bbaas/pkg/request"

	"github.com/antchfx/htmlquery"
)

var scope []string
var s string

// Scrape tries to grab scope table for a given program on intigriti.com
func Scrape(url string) string {
	re := regexp.MustCompile(`([\w-]+)\/([\w-]+)\/([\w-]+$)`)
	match := re.FindStringSubmatch(url)
	business := strings.ToLower(match[1])
	program := strings.ToLower(match[2])
	endpoint := "https://app.intigriti.com/programs/" + business + "/" + program + "/" + "detail"

	// GET request to endpoint
	resp, status := req.GET(endpoint)

	// check bad status code
	if status != 200 {
		errors.BadStatusCode(url, status)
	}

	// parse response body to xQuery doc
	doc, _ := htmlquery.Parse(strings.NewReader(resp))

	// xQuery to grab in-scope and out-of-scope tables
	blob := htmlquery.Find(doc, "//div[@class='domains']")

	for _, t := range blob {
		s = s + htmlquery.InnerText(t)
	}

	s = strings.Replace(s, "\n", " ", -1)
	re1 := regexp.MustCompile(`(.*)`)
	re2 := regexp.MustCompile(`(Out of Scope(.*))`)

	inscope := re1.FindString(s)
	outscope := re2.FindString(s)

	inscope = re2.ReplaceAllString(inscope, "$3") //remove out-of-scope items
	scope = append(scope, "!INCLUDE")
	scope = append(scope, inscope)
	scope = append(scope, "!EXCLUDE")
	scope = append(scope, outscope)

	return strings.Join(scope, "\n")
}
