package sprig

import "html/template"

// unescape the path
func toHtml(path string) interface{} {
	return template.HTML(path)
}
