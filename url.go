package sprig

import (
	"fmt"
	"net/url"
	"reflect"
)

func dictGetOrEmpty(dict map[string]interface{}, key string) string {
	value, ok := dict[key]
	if !ok {
		return ""
	}
	tp := reflect.TypeOf(value).Kind()
	if tp != reflect.String {
		panic(fmt.Sprintf("failed to parse key %q: expected string type, got %s", key, tp.String()))
	}
	return reflect.ValueOf(value).String()
}

// parses given URL to return dict object
func urlParse(v string) map[string]interface{} {
	dict := map[string]interface{}{}
	parsedURL, err := url.Parse(v)
	if err != nil {
		panic(fmt.Sprintf("failed to parse URL: %v", err))
	}
	dict["scheme"] = parsedURL.Scheme
	dict["host"] = parsedURL.Host
	dict["hostname"] = parsedURL.Hostname()
	dict["path"] = parsedURL.Path
	dict["query"] = parsedURL.RawQuery
	dict["opaque"] = parsedURL.Opaque
	dict["fragment"] = parsedURL.Fragment
	if parsedURL.User != nil {
		dict["userinfo"] = parsedURL.User.String()
	} else {
		dict["userinfo"] = ""
	}

	return dict
}

// join given dict to URL string
func urlJoin(d map[string]interface{}) string {
	resURL := url.URL{
		Scheme:   dictGetOrEmpty(d, "scheme"),
		Host:     dictGetOrEmpty(d, "host"),
		Path:     dictGetOrEmpty(d, "path"),
		RawQuery: dictGetOrEmpty(d, "query"),
		Opaque:   dictGetOrEmpty(d, "opaque"),
		Fragment: dictGetOrEmpty(d, "fragment"),
	}
	userinfo := dictGetOrEmpty(d, "userinfo")
	var user *url.Userinfo
	if userinfo != "" {
		tempURL, err := url.Parse(fmt.Sprintf("proto://%s@host", userinfo))
		if err != nil {
			panic(fmt.Sprintf("failed to parse userinfo in dict: %v", err))
		}
		user = tempURL.User
	}

	resURL.User = user
	return resURL.String()
}

// path escapes the strings
func pathEscape(s string) string {
	return url.PathEscape(s)
}

// path unescapes the strings
func pathUnescape(s string) string {
	ue, err := url.PathUnescape(s)
	if err != nil {
		return s
	}

	return ue
}

// urlEncode encodes special characters in URL query parameters while preserving URL structure.
// This is useful for Markdown links where unencoded spaces/quotes in URLs break link parsing.
func urlEncode(rawURL string) string {
	if rawURL == "" {
		return ""
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return url.QueryEscape(rawURL)
	}

	// Re-encode query parameters properly
	if parsed.RawQuery != "" {
		values, err := url.ParseQuery(parsed.RawQuery)
		if err == nil {
			parsed.RawQuery = values.Encode()
		}
	}

	return parsed.String()
}

// urlDecode decodes percent-encoded characters in URL query parameters and path.
func urlDecode(encodedURL string) string {
	if encodedURL == "" {
		return ""
	}

	parsed, err := url.Parse(encodedURL)
	if err != nil {
		decoded, err := url.QueryUnescape(encodedURL)
		if err != nil {
			return encodedURL
		}
		return decoded
	}

	// Decode path
	if parsed.Path != "" {
		decodedPath, err := url.PathUnescape(parsed.Path)
		if err == nil {
			parsed.Path = decodedPath
		}
	}

	// Decode query parameters
	if parsed.RawQuery != "" {
		decoded, err := url.QueryUnescape(parsed.RawQuery)
		if err == nil {
			parsed.RawQuery = decoded
		}
	}

	return parsed.String()
}
