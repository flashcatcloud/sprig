package sprig

import (
	"fmt"
	"html/template"
)

// toHtml converts the first non-empty value to template.HTML, bypassing HTML escaping.
func toHtml(vals ...interface{}) interface{} {
	for _, v := range vals {
		if v == nil {
			continue
		}
		if v == "" {
			continue
		}
		return template.HTML(fmt.Sprintf("%v", v))
	}

	return ""
}
