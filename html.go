package sprig

import (
	"fmt"
	"html/template"
)

// unescape the vals
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
