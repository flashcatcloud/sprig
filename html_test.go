package sprig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToHtml(t *testing.T) {
	// testing that function is exported and working properly
	assert.NoError(t, runt(
		`{{ toHtml "CPU idle &gt; 90%" }}`,
		"CPU idle > 90%"))

	// testing scenarios
	for url, expected := range urlTests {
		assert.EqualValues(t, expected, urlParse(url))
	}
}
