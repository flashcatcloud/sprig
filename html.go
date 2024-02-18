package sprig

import "html/template"

// unescape the vals
func toHtml(vals ...string) interface{} {
	for _, v := range vals {
		if v != "" {
			return template.HTML(v)
		}
	}
	return ""
}
