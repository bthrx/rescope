//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2019 root4loot
//

package bugcrowd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"

	errors "github.com/root4loot/rescope/internal/bbaas/pkg/errors"
	req "github.com/root4loot/rescope/internal/bbaas/pkg/request"
)

// Scrape returns a string containing scope that was scraped from the given program on bugcrowd.com
func Scrape(url string) string {
	re := regexp.MustCompile(`([\w-]+)\/([\w-]+$)`)
	match := re.FindStringSubmatch(url)
	program := match[2]
	endpoint := "https://bugcrowd.com/" + program
	var scope []string

	// GET request to endpoint
	resp, status := req.GET(endpoint)

	// check bad status code
	if status != 200 {
		errors.BadStatusCode(url, status)
	}

	// parse response body to xQuery doc
	doc, _ := htmlquery.Parse(strings.NewReader(resp))

	// xQuery to grab in-scope and out-of-scope tables
	blob := htmlquery.Find(doc, "//div[@data-react-class='ResearcherTargetGroups']")

	for _, item := range blob {
		s := fmt.Sprintf("a %s", item)
		s = strings.Replace(s, "\\u003c", "<", -1)
		s = strings.Replace(s, "\\u003e", ">", -1)

		// remove unwanted tags from blob
		re1 := regexp.MustCompile(`({"tags":\[(.*?)})`)
		re2 := regexp.MustCompile(`(category":"other",(.*?)})`)
		re3 := regexp.MustCompile(`("uri":null,"target":{"id"(.*?)})`)
		s = re1.ReplaceAllString(s, "$2")
		s = re2.ReplaceAllString(s, "$2")
		s = re3.ReplaceAllString(s, "$3")

		re4 := regexp.MustCompile(`"in_scope":true(.*)`)
		re5 := regexp.MustCompile(`"(in_scope":false(.*))`)
		inscope := re4.FindAllString(s, -1)
		outscope := re5.FindAllString(s, -1)

		scope = append(scope, "!INCLUDE")
		for _, item := range inscope {
			item = re5.ReplaceAllString(item, "$3") //remove out-of-scope items
			scope = append(scope, item)
		}

		scope = append(scope, "!EXCLUDE")
		for _, item := range outscope {
			scope = append(scope, item)
		}
	}

	return strings.Join(scope, "\n")
}
