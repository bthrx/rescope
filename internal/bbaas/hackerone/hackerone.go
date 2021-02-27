//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2021 root4loot
//

package hackerone

import (
	"reflect"
	"regexp"
	"strings"
	"unsafe"

	req "github.com/root4loot/rescope/internal/bbaas/pkg/request"
)

var scope []string
var isAssetURL bool

// Scrape tries to grab scope table for a given program on hackerone.com
func Scrape(url string) string {
	re := regexp.MustCompile(`([\w-]+)\/([\w-]+$)`)
	match := re.FindStringSubmatch(url)
	program := strings.ToLower(match[2])
	endpoint := "https://hackerone.com/graphql?"

	// clear global slice
	scope = nil

	var data = []byte(`{  
		"query":"query Team_assets($first_0:Int!) {query {id,...F0}} fragment F0 on Query {_teamAgUhl:team(handle:\"` + program + `\") {handle,_structured_scope_versions2ZWKHQ:structured_scope_versions(archived:false) {max_updated_at},_structured_scopeszxYtW:structured_scopes(first:$first_0,archived:false,eligible_for_submission:true) {edges {node {asset_type, asset_identifier}},pageInfo {hasNextPage,hasPreviousPage}},_structured_scopes3FF98f:structured_scopes(first:$first_0,archived:false,eligible_for_submission:false) {edges {node {asset_type,asset_identifier,},},},},}",
		"variables":{  
		   "first_0":1337
		}
	 }`)

	resB, _ := (req.POST(endpoint, data))
	resS := BytesToString(resB)

	re = regexp.MustCompile(`\"edges":\[(.*?)\]`)
	scopeSplit := re.FindAllString(resS, -1)
	re = regexp.MustCompile(`asset_type":"(URL|CIDR|IP|IP-RANGE|RANGE)","asset_identifier":"(.*?)"`)

	inScope := re.FindAllStringSubmatch(scopeSplit[0], -1)
	outScope := re.FindAllStringSubmatch(scopeSplit[1], -1)

	scope = append(scope, "!INCLUDE")
	for _, m := range inScope {
		scope = append(scope, m[2])
	}

	scope = append(scope, "!EXCLUDE")
	for _, m := range outScope {
		scope = append(scope, m[2])
	}

	return strings.Join(scope, "\n")
}

// BytesToString converts byte array to string
func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}
