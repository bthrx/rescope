//
// Written By : Daniel Antonsen (@root4loot)
//
// Distributed Under MIT License
// Copyrights (C) 2021 root4loot
//

package yeswehack

import (
	"testing"
)

// TestHackerone scrape
func TestYeswehack(t *testing.T) {
	Scrape("https://yeswehack.com/programs/yes-we-hack")
}
