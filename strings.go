package sprig

import (
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	util "github.com/Masterminds/goutils"
)

func base64encode(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func base64decode(v string) string {
	data, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func base32encode(v string) string {
	return base32.StdEncoding.EncodeToString([]byte(v))
}

func base32decode(v string) string {
	data, err := base32.StdEncoding.DecodeString(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func abbrev(width int, s string) string {
	if width < 4 {
		return s
	}
	r, _ := util.Abbreviate(s, width)
	return r
}

func abbrevboth(left, right int, s string) string {
	if right < 4 || left > 0 && right < 7 {
		return s
	}
	r, _ := util.AbbreviateFull(s, left, right)
	return r
}

func initials(s string) string {
	// Wrap this just to eliminate the var args, which templates don't do well.
	return util.Initials(s)
}

func randAlphaNumeric(count int) string {
	// It is not possible, it appears, to actually generate an error here.
	r, _ := util.CryptoRandomAlphaNumeric(count)
	return r
}

func randAlpha(count int) string {
	r, _ := util.CryptoRandomAlphabetic(count)
	return r
}

func randAscii(count int) string {
	r, _ := util.CryptoRandomAscii(count)
	return r
}

func randNumeric(count int) string {
	r, _ := util.CryptoRandomNumeric(count)
	return r
}

func untitle(str string) string {
	return util.Uncapitalize(str)
}

func quote(str ...interface{}) string {
	out := make([]string, 0, len(str))
	for _, s := range str {
		if s != nil {
			out = append(out, fmt.Sprintf("%q", strval(s)))
		}
	}
	return strings.Join(out, " ")
}

func squote(str ...interface{}) string {
	out := make([]string, 0, len(str))
	for _, s := range str {
		if s != nil {
			out = append(out, fmt.Sprintf("'%v'", s))
		}
	}
	return strings.Join(out, " ")
}

func cat(v ...interface{}) string {
	v = removeNilElements(v)
	r := strings.TrimSpace(strings.Repeat("%v ", len(v)))
	return fmt.Sprintf(r, v...)
}

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

func nindent(spaces int, v string) string {
	return "\n" + indent(spaces, v)
}

func replace(old, new, src string) string {
	return strings.Replace(src, old, new, -1)
}

func plural(one, many string, count int) string {
	if count == 1 {
		return one
	}
	return many
}

func strslice(v interface{}) []string {
	switch v := v.(type) {
	case []string:
		return v
	case []interface{}:
		b := make([]string, 0, len(v))
		for _, s := range v {
			if s != nil {
				b = append(b, strval(s))
			}
		}
		return b
	default:
		val := reflect.ValueOf(v)
		switch val.Kind() {
		case reflect.Array, reflect.Slice:
			l := val.Len()
			b := make([]string, 0, l)
			for i := 0; i < l; i++ {
				value := val.Index(i).Interface()
				if value != nil {
					b = append(b, strval(value))
				}
			}
			return b
		default:
			if v == nil {
				return []string{}
			}

			return []string{strval(v)}
		}
	}
}

func removeNilElements(v []interface{}) []interface{} {
	newSlice := make([]interface{}, 0, len(v))
	for _, i := range v {
		if i != nil {
			newSlice = append(newSlice, i)
		}
	}
	return newSlice
}

func strval(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func trunc(c int, s string) string {
	if c < 0 && len(s)+c > 0 {
		return s[len(s)+c:]
	}
	if c >= 0 && len(s) > c {
		return s[:c]
	}
	return s
}

func join(sep string, v interface{}) string {
	return strings.Join(strslice(v), sep)
}

func split(sep, orig string) map[string]string {
	parts := strings.Split(orig, sep)
	res := make(map[string]string, len(parts))
	for i, v := range parts {
		res["_"+strconv.Itoa(i)] = v
	}
	return res
}

func splitn(sep string, n int, orig string) map[string]string {
	parts := strings.SplitN(orig, sep, n)
	res := make(map[string]string, len(parts))
	for i, v := range parts {
		res["_"+strconv.Itoa(i)] = v
	}
	return res
}

// substring creates a substring of the given string.
//
// If start is < 0, this calls string[:end].
//
// If start is >= 0 and end < 0 or end bigger than s length, this calls string[start:]
//
// Otherwise, this calls string[start, end].
func substring(start, end int, s string) string {
	length := len(s)

	// Handle both negative case: return full string
	if start < 0 && end < 0 {
		return s
	}

	// Normalize negative indices first
	normalizedStart := start
	normalizedEnd := end
	if normalizedStart < 0 {
		normalizedStart = 0
	}
	if normalizedEnd < 0 {
		normalizedEnd = length
	}

	// Handle start > end case: return from 0 to start
	// Only check this after normalizing negative indices
	if normalizedStart > normalizedEnd {
		normalizedStart, normalizedEnd = 0, normalizedStart
	}

	// Check bounds
	if normalizedStart > length {
		return ""
	}
	if normalizedEnd > length {
		normalizedEnd = length
	}
	if normalizedStart >= normalizedEnd {
		return ""
	}

	// Ensure we don't cut UTF-8 characters in the middle
	// Adjust start forward if it's in the middle of a UTF-8 character
	if normalizedStart > 0 && normalizedStart < length {
		// Check if start is at a UTF-8 character boundary
		if !utf8.RuneStart(s[normalizedStart]) {
			// Find the start of the next valid UTF-8 character
			for normalizedStart < length && normalizedStart < normalizedEnd && !utf8.RuneStart(s[normalizedStart]) {
				normalizedStart++
			}
			if normalizedStart >= normalizedEnd {
				return ""
			}
		}
	}

	// Adjust end if it's in the middle of a UTF-8 character
	if normalizedEnd > normalizedStart && normalizedEnd < length {
		// Check if end is at a UTF-8 character boundary
		if !utf8.RuneStart(s[normalizedEnd]) {
			// Find the start of the current character by going backward
			tempEnd := normalizedEnd
			for tempEnd > normalizedStart && !utf8.RuneStart(s[tempEnd]) {
				tempEnd--
			}
			// Now find the end of this character
			if tempEnd >= normalizedStart && tempEnd < length {
				_, size := utf8.DecodeRuneInString(s[tempEnd:])
				normalizedEnd = tempEnd + size
				if normalizedEnd > length {
					normalizedEnd = length
				}
			} else {
				// If we can't find a valid character start, just use the original end
				normalizedEnd = length
			}
		}
	}

	// Extract substring after UTF-8 boundary adjustments
	if normalizedStart >= normalizedEnd || normalizedStart > length {
		return ""
	}
	if normalizedEnd > length {
		normalizedEnd = length
	}

	return s[normalizedStart:normalizedEnd]
}

// substrRune creates a substring of the given string based on runes.
//
// If start is < 0, this calls string[:end].
//
// If start is >= 0 and end < 0 or end bigger than s length, this calls string[start:]
//
// Otherwise, this calls string[start, end].
//
// This is a multi-byte safe version of substring.
func substrRune(start, end int, s string) string {
	runes := []rune(s)
	l := len(runes)
	if start < 0 {
		if end > l {
			end = l
		}
		if end < 0 {
			return ""
		}
		return string(runes[:end])
	}
	if end < 0 || end > l {
		if start > l {
			return ""
		}
		return string(runes[start:])
	}
	if start > end {
		return ""
	}
	return string(runes[start:end])
}
