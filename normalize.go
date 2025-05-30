package sprig

import (
	"regexp"
	"strings"
)

// Define regular expressions based on pattern.md
var (
	// URL (simplified, catches http/https prefix and common chars)
	// Needs careful balance to avoid matching too much or too little.
	urlRegex = regexp.MustCompile(`https?://[a-zA-Z0-9\-\._~:/?#[\]@!$&'()*+,;=%]+`)

	// IP Address v4/v6 (basic) - Combine with port handling
	ipRegex = regexp.MustCompile(`(\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b|\b(?:[a-fA-F0-9]{1,4}:){2,}[a-fA-F0-9:/%.]{1,}\b)(:\d+)?`) // Improved IPv6 matching slightly

	// Matches common source file paths like /path/to/file.js, ./src/component.tsx, C:\\path\\file.vue, etc.
	sourceFilePathRegex = regexp.MustCompile(`(?:[a-zA-Z]:)?(?:[/\\][\w.-]+)+?\.(?:jsx?|tsx?|vue|svelte|css|scss|less|sass|html?|json|svg|graphql|gql|map)\b`)

	// File paths or names - Tries to catch common structures like /path/file.ext, file-hash.ext, file.hash.ext
	// This is inherently difficult and might need refinement based on actual data.
	// Prioritize replacing full path-like structures found.
	// Adjusted to be less greedy and focus on structures with slashes or extensions/hashes.
	filePathRegex = regexp.MustCompile(`(?:[/\.\w\-]+)*([\w\-]+(?:[-.][0-9a-f]{6,}|\.[a-zA-Z0-9]{1,5}))(?:[:\s]|$)`)

	// UUID format regex
	uuidRegex = regexp.MustCompile(`\b[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}\b`)

	// MongoDB ObjectID format (24 hex characters)
	objectIDRegex = regexp.MustCompile(`\b[0-9a-f]{24}\b`)

	// Pure numeric segment regex (long standalone numbers)
	numberRegex = regexp.MustCompile(`\b\d+\b`)

	// Hash values regex (MD5, SHA1, SHA256, etc.) - includes common lengths
	hashRegex = regexp.MustCompile(`\b[0-9a-f]{32}\b|\b[0-9a-f]{40}\b|\b[0-9a-f]{64}\b`)

	// Trace/Span ID format (typically 16 or 32 hex chars)
	traceIDRegex = regexp.MustCompile(`\b[0-9a-f]{16}\b|\b[0-9a-f]{32}\b`) // More specific lengths

	// Hexadecimal segment with length >= 6 (likely an identifier, avoid replacing hashes already caught)
	hexLongSegmentRegex = regexp.MustCompile(`\b[0-9a-f]{6,}\b`)

	// Refined to avoid matching colons directly after quotes/spaces (common in JSON).
	// Relies on URL/IP regex running first to handle host:port cases.
	lineNumRegex = regexp.MustCompile(`(?:[^\"'\s:]:|[@]|at\s|line\s|#)(\d+)`) // Avoid colon after quote/space/colon

	// Regex to clean multiple spaces
	multiSpaceRegex = regexp.MustCompile(`\s{2,}`)

	// Regex to find the start of a common stack trace format
	stackTraceStartPattern = "\n    at "
)

// Placeholders
const (
	urlPlaceholder         = "{URL}"
	ipPlaceholder          = "{IP}"
	filePlaceholder        = "{FILE}"
	hashPlaceholder        = "{HASH}"   // Generic placeholder for hashes, long numbers, IDs
	linePlaceholder        = "{LINE}"   // Placeholder for line number info
	numberPlaceholder      = "{NUMBER}" // Placeholder for numbers (typically long ones)
	stackFramesPlaceholder = "{StackFrames}"
)

// normalizeMessage processes an error message to extract its stable part using regex.
func normalizeMessage(message string) string {
	if message == "" {
		return ""
	}

	normalized := message

	// --- Apply existing normalization rules first ---

	// 1. URLs
	normalized = urlRegex.ReplaceAllString(normalized, urlPlaceholder)

	// 2. IPs (including port)
	normalized = ipRegex.ReplaceAllString(normalized, ipPlaceholder)

	// 3. Source File Paths (Specific extensions like .js, .tsx, .vue, etc.) - Run before generic file paths
	normalized = sourceFilePathRegex.ReplaceAllString(normalized, filePlaceholder)

	// 4. Generic File Paths / Names (often contain other patterns like hashes/numbers)
	normalized = filePathRegex.ReplaceAllString(normalized, filePlaceholder)

	// 5. Line number indicators (like :123, @123, at 123, line 123, #123)
	// Replace the line number part, keeping the indicator if possible (more complex, currently replaces match)
	normalized = lineNumRegex.ReplaceAllString(normalized, linePlaceholder) // Simple replacement for now

	// 6. Standalone long numbers
	normalized = numberRegex.ReplaceAllString(normalized, numberPlaceholder)

	// 7. Specific ID formats (UUID, ObjectID, TraceID)
	normalized = uuidRegex.ReplaceAllString(normalized, hashPlaceholder)
	normalized = objectIDRegex.ReplaceAllString(normalized, hashPlaceholder)
	normalized = traceIDRegex.ReplaceAllString(normalized, hashPlaceholder) // Place before general hash/hex

	// 8. Common hash lengths (MD5, SHA1, SHA256)
	normalized = hashRegex.ReplaceAllString(normalized, hashPlaceholder)

	// 9. Generic long hexadecimal strings (that weren't matched as specific IDs/hashes)
	normalized = hexLongSegmentRegex.ReplaceAllString(normalized, hashPlaceholder)

	// 10. Stack frames
	stackStartIndex := strings.Index(normalized, stackTraceStartPattern)
	if stackStartIndex != -1 {
		// Keep the part before the stack trace starts
		normalized = normalized[:stackStartIndex]
		// Append the placeholder
		normalized += " " + stackFramesPlaceholder // Add a space before placeholder
	}

	// --- Final cleanup ---

	// Replace multiple spaces with a single space
	normalized = multiSpaceRegex.ReplaceAllString(normalized, " ")

	// Trim leading/trailing spaces
	normalized = strings.TrimSpace(normalized)

	return normalized
}
