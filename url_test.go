package sprig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var urlTests = map[string]map[string]interface{}{
	"proto://auth@host:80/path?query#fragment": {
		"fragment": "fragment",
		"host":     "host:80",
		"hostname": "host",
		"opaque":   "",
		"path":     "/path",
		"query":    "query",
		"scheme":   "proto",
		"userinfo": "auth",
	},
	"proto://host:80/path": {
		"fragment": "",
		"host":     "host:80",
		"hostname": "host",
		"opaque":   "",
		"path":     "/path",
		"query":    "",
		"scheme":   "proto",
		"userinfo": "",
	},
	"something": {
		"fragment": "",
		"host":     "",
		"hostname": "",
		"opaque":   "",
		"path":     "something",
		"query":    "",
		"scheme":   "",
		"userinfo": "",
	},
	"proto://user:passwor%20d@host:80/path": {
		"fragment": "",
		"host":     "host:80",
		"hostname": "host",
		"opaque":   "",
		"path":     "/path",
		"query":    "",
		"scheme":   "proto",
		"userinfo": "user:passwor%20d",
	},
	"proto://host:80/pa%20th?key=val%20ue": {
		"fragment": "",
		"host":     "host:80",
		"hostname": "host",
		"opaque":   "",
		"path":     "/pa th",
		"query":    "key=val%20ue",
		"scheme":   "proto",
		"userinfo": "",
	},
}

func TestUrlParse(t *testing.T) {
	// testing that function is exported and working properly
	assert.NoError(t, runt(
		`{{ index ( urlParse "proto://auth@host:80/path?query#fragment" ) "host" }}`,
		"host:80"))

	// testing scenarios
	for url, expected := range urlTests {
		assert.EqualValues(t, expected, urlParse(url))
	}
}

func TestUrlJoin(t *testing.T) {
	tests := map[string]string{
		`{{ urlJoin (dict "fragment" "fragment" "host" "host:80" "path" "/path" "query" "query" "scheme" "proto") }}`:       "proto://host:80/path?query#fragment",
		`{{ urlJoin (dict "fragment" "fragment" "host" "host:80" "path" "/path" "scheme" "proto" "userinfo" "ASDJKJSD") }}`: "proto://ASDJKJSD@host:80/path#fragment",
	}
	for tpl, expected := range tests {
		assert.NoError(t, runt(tpl, expected))
	}

	for expected, urlMap := range urlTests {
		assert.EqualValues(t, expected, urlJoin(urlMap))
	}
}

func TestPathEscape(t *testing.T) {
	// testing that function is exported and working properly
	assert.NoError(t, runt(
		`{{ pathEscape "CPU idle > 90%" }}`,
		"CPU%20idle%20%3E%2090%25"))

	// testing scenarios
	for url, expected := range urlTests {
		assert.EqualValues(t, expected, urlParse(url))
	}
}

func TestPathUnEscape(t *testing.T) {
	// testing that function is exported and working properly
	assert.NoError(t, runt(
		`{{ pathUnescape "CPU%20idle%20%3E%2090%25" }}`,
		"CPU idle > 90%"))

	// testing scenarios
	for url, expected := range urlTests {
		assert.EqualValues(t, expected, urlParse(url))
	}
}

func TestUrlEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple url without special chars",
			input:    "https://example.com/path?key=value",
			expected: "https://example.com/path?key=value",
		},
		{
			name:     "url with spaces in query",
			input:    "https://example.com/path?query=hello world",
			expected: "https://example.com/path?query=hello+world",
		},
		{
			name:     "url with quotes in query",
			input:    `https://example.com/path?query=type: "*"`,
			expected: "https://example.com/path?query=type%3A+%22%2A%22",
		},
		{
			name:     "url with parentheses in query",
			input:    "https://example.com/path?query=(a OR b)",
			expected: "https://example.com/path?query=%28a+OR+b%29",
		},
		{
			name:     "complex lucene query url",
			input:    `https://flashcat.net/log/explorer?data_source_id=36&data_source_name=elasticsearch&end=now&index_pattern=458&mode=index_pattern&query=client_type: "*" AND NOT client_type: "WX_MINI"  AND log_type: "RAW_RESPONSE"  AND NOT (response_code:200 OR response_code:100) AND (   url_path: "/mop/v2/store/list"   OR url_path: "/mod/v4/store/list"   OR url_path: "/mop/v1/store/data_by_store_ids"     OR url_path: "/mod/v1/store/detail" )&start=now-1h&syntax=lucene&__execute__=true`,
			expected: "https://flashcat.net/log/explorer?__execute__=true&data_source_id=36&data_source_name=elasticsearch&end=now&index_pattern=458&mode=index_pattern&query=client_type%3A+%22%2A%22+AND+NOT+client_type%3A+%22WX_MINI%22++AND+log_type%3A+%22RAW_RESPONSE%22++AND+NOT+%28response_code%3A200+OR+response_code%3A100%29+AND+%28+++url_path%3A+%22%2Fmop%2Fv2%2Fstore%2Flist%22+++OR+url_path%3A+%22%2Fmod%2Fv4%2Fstore%2Flist%22+++OR+url_path%3A+%22%2Fmop%2Fv1%2Fstore%2Fdata_by_store_ids%22+++++OR+url_path%3A+%22%2Fmod%2Fv1%2Fstore%2Fdetail%22+%29&start=now-1h&syntax=lucene",
		},
		{
			name:     "url with fragment",
			input:    "https://example.com/path?query=test#section",
			expected: "https://example.com/path?query=test#section",
		},
		{
			name:     "url without query string",
			input:    "https://example.com/path",
			expected: "https://example.com/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := urlEncode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}

	// testing that function is exported and working properly via template
	assert.NoError(t, runt(
		`{{ urlEncode "https://example.com/path?query=hello world" }}`,
		"https://example.com/path?query=hello+world"))
}

func TestUrlDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple url without encoding",
			input:    "https://example.com/path?key=value",
			expected: "https://example.com/path?key=value",
		},
		{
			name:     "url with plus sign in query",
			input:    "https://example.com/path?query=hello+world",
			expected: "https://example.com/path?query=hello world",
		},
		{
			name:     "url with percent encoding in query",
			input:    "https://example.com/path?query=type%3A+%22%2A%22",
			expected: `https://example.com/path?query=type: "*"`,
		},
		{
			name:     "url with encoded parentheses",
			input:    "https://example.com/path?query=%28a+OR+b%29",
			expected: "https://example.com/path?query=(a OR b)",
		},
		{
			name:     "url with fragment",
			input:    "https://example.com/path?query=test#section",
			expected: "https://example.com/path?query=test#section",
		},
		{
			name:     "url without query string",
			input:    "https://example.com/path",
			expected: "https://example.com/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := urlDecode(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}

	// testing that function is exported and working properly via template
	assert.NoError(t, runt(
		`{{ urlDecode "https://example.com/path?query=hello+world" }}`,
		"https://example.com/path?query=hello world"))
}
