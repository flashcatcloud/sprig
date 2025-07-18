package sprig

import (
	"regexp"
	"strings"
)

// Define regular expressions with better precision and ordering
var (
	// URL (simplified, catches http/https prefix and common chars)
	urlRegex = regexp.MustCompile(`https?://[a-zA-Z0-9\-\._~:/?#[\]@!$&'()*+,;=%]+`)

	// IP Address v4/v6 with optional port (avoid matching time formats like 10:16:04)
	ipv4Regex = regexp.MustCompile(`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b(?::\d+)?`)
	ipv6Regex = regexp.MustCompile(`\b(?:[a-fA-F0-9]{1,4}:){2,}[a-fA-F0-9:/%.]*[a-fA-F][a-fA-F0-9:/%.]*\b|\b(?:[a-fA-F0-9]{1,4}:){7}[a-fA-F0-9]{1,4}\b`)

	// File paths with specific extensions (including Java stack trace format)
	sourceFilePathRegex = regexp.MustCompile(`(?:[a-zA-Z]:)?(?:[/\\][\w.-]+)+?\.(?:jsx?|tsx?|vue|svelte|css|scss|less|sass|html?|json|svg|graphql|gql|map|java|py|rb|php|go|rs|kt|swift|cpp|cc|c|h|hpp)\b|\b\w+\.(?:java|py|rb|php|go|rs|kt|swift|cpp|cc|c|h|hpp)\b`)

	// UUID format (36 characters with dashes)
	uuidRegex = regexp.MustCompile(`\b[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}\b`)

	// JWT token format (header.payload.signature, each part is base64url encoded)
	jwtRegex = regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\b`)

	// MongoDB ObjectID format (24 hex characters)
	objectIDRegex = regexp.MustCompile(`\b[0-9a-f]{24}\b`)

	// Trace/Span ID format (32 hex chars)
	traceIDRegex = regexp.MustCompile(`\b[0-9a-f]{32}\b`)

	// Hash values (SHA1: 40, SHA256: 64 hex chars)
	hashRegex = regexp.MustCompile(`\b[0-9a-f]{40}\b|\b[0-9a-f]{64}\b`)

	// Long hexadecimal strings (8+ chars, not caught by specific patterns above)
	hexLongRegex = regexp.MustCompile(`\b[0-9a-f]{8,}\b`)

	// Long identifiers with underscores (like default_26df169d0e0c35325d5d80fb1f3beda3_dc4a0677b18b4c8e985f63960647fe34)
	longIdentifierRegex = regexp.MustCompile(`\b[a-zA-Z0-9_]{40,}\b`)

	// API keys and tokens - more specific pattern for uppercase alphanumeric strings
	// Only match strings that are mostly uppercase and contain at least some numbers
	apiKeyRegex = regexp.MustCompile(`\b[A-Z]{2,}[A-Z0-9]{6,}[0-9][A-Z0-9]*\b`)

	// Base64 encoded strings (more specific: 20+ chars, containing +/, ending with =)
	base64Regex = regexp.MustCompile(`\b[A-Za-z0-9+/]{20,}(?:[+/][A-Za-z0-9+/]*)?={1,2}\b`)

	// Task IDs and similar alphanumeric identifiers (like taskId:528jLlkO2sr7OjGV, abc123)
	// Match strings that have mixed case and numbers, at least 6 chars, avoid pure hex
	taskIdRegex = regexp.MustCompile(`\b[a-zA-Z]*[0-9]+[a-zA-Z]+[a-zA-Z0-9]{4,}\b|\b[a-zA-Z]+[0-9]+[a-zA-Z0-9]{2,}\b`)

	// Date and time patterns
	// ISO 8601 format: 2024-04-07T19:00:05, 2024-12-17T17:56:28
	iso8601DateRegex = regexp.MustCompile(`\b\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d{3})?(?:Z|[+-]\d{2}:\d{2})?\b`)
	// Traditional format: Fri Jul 18 10:16:04 HKT 2025
	traditionalDateRegex = regexp.MustCompile(`(?:Mon|Tue|Wed|Thu|Fri|Sat|Sun)\s+(?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)\s+\d{1,2}\s+\d{2}:\d{2}:\d{2}\s+[A-Z]{3,4}\s+\d{4}`)
	// RFC2822 format: Mon, 02 Jan 2006 15:04:05 MST
	rfc2822DateRegex = regexp.MustCompile(`(?:Mon|Tue|Wed|Thu|Fri|Sat|Sun),\s+\d{1,2}\s+(?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)\s+\d{4}\s+\d{2}:\d{2}:\d{2}\s+[A-Z]{3,4}`)
	// Complete datetime formats: 2024/04/07 14:30:25.123, 2024-04-07 14:30:25
	completeDateTimeRegex = regexp.MustCompile(`\b\d{4}[-/]\d{1,2}[-/]\d{1,2}\s+\d{1,2}:\d{2}:\d{2}(?:\.\d{1,3})?\b`)
	// Date only formats: 2024-04-07, 2024/04/07, 04/07/2024, 07-04-2024
	dateOnlyRegex = regexp.MustCompile(`\b\d{4}[-/]\d{1,2}[-/]\d{1,2}\b|\b\d{1,2}[-/]\d{1,2}[-/]\d{4}\b`)
	// Month name formats: Apr 7, 2024, 7 Apr 2024, April 7, 2024, 7 April 2024
	monthNameRegex = regexp.MustCompile(`\b(?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec|January|February|March|April|May|June|July|August|September|October|November|December)\s+\d{1,2},?\s+\d{4}\b|\b\d{1,2}\s+(?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec|January|February|March|April|May|June|July|August|September|October|November|December)\s+\d{4}\b`)
	// Time only formats: 19:00:05, 19:00:05.123, 7:00:05 PM, 7:00 PM
	timeOnlyRegex = regexp.MustCompile(`\b\d{1,2}:\d{2}:\d{2}(?:\.\d{1,3})?\b|\b\d{1,2}:\d{2}(?::\d{2})?\s*(?:AM|PM)\b`)
	// Chinese date format: 2024年04月07日, 2024年4月7日
	chineseDateRegex = regexp.MustCompile(`\d{4}年\d{1,2}月\d{1,2}日`)

	// Domain names (be more conservative - only match common TLDs and include cc)
	domainRegex = regexp.MustCompile(`\b[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)*(?:com|ai|top|dev|sh|org|net|edu|gov|mil|int|cloud|io|co|cn|jp|de|uk|fr|au|ca|ru|it|br|in|mx|kr|sg|hk|tw|th|my|id|ph|vn|bd|pk|ng|za|eg|ma|dz|ke|tz|ug|rw|mw|zm|zw|bw|sz|ls|mz|mg|mu|sc|km|dj|so|et|er|sd|ss|ly|tn|mr|ml|ne|td|cf|cm|ga|gq|st|cv|gw|gm|sl|lr|ci|gh|tg|bj|bf|sn|gn|ge|am|az|cc|us|es|pl|nl|be|at|ch|se|no|dk|fi|pt|gr|cz|hu|ro|bg|hr|si|sk|ee|lv|lt|ie|is|mt|cy|lu|tr|il|ae|sa|kw|qa|bh|om|jo|lb|sy|iq|ir|af|np|bt|lk|mv|mm|la|kh|bn|tl|pg|fj|vu|ws|to|tv|ki|nr|pw|fm|mh|as|gu|mp|vi|pr|ck|nu|tk|pn|info|biz|name|pro|museum|aero|coop|jobs|travel|mobi|tel|xxx|asia|cat|post|app|blog|site|online|store|tech|digital|media|news|live|video|music|games|social|email|web|world|global|city|space|life|work|business|company|corp|inc|ltd|group|team|club|community|network|solutions|services|consulting|agency|studio|design|creative|art|photo|gallery|shop|market|trade|finance|money|bank|investment|insurance|real|estate|property|hotel|restaurant|food|health|medical|care|fitness|sport|auto|car|bike|fashion|beauty|style|luxury|premium|vip|gold|diamond|jewelry|watch|home|house|garden|kids|family|baby|pet|love|dating|wedding|party|event|education|school|university|college|academy|training|course|book|library|science|research|technology|software|hardware|computer|internet|hosting|mobile|phone|game|download|free|cheap|sale|deal|discount|coupon|gift|win|lottery|casino|bet)\b`)

	// Stack trace start pattern
	stackTraceStartPattern = "\n    at "

	// Special file patterns that should be empty
	specialFilePattern = regexp.MustCompile(`^/t\d+/[\w-]+/\?/\?/\?\.\w+$`)

	// Regex to clean multiple spaces
	multiSpaceRegex = regexp.MustCompile(`\s{2,}`)
)

// Placeholders
const (
	urlPlaceholder         = "{URL}"
	ipPlaceholder          = "{IP}"
	domainPlaceholder      = "{DOMAIN}"
	filePlaceholder        = "{FILE}"
	hashPlaceholder        = "{HASH}"
	numberPlaceholder      = "{NUMBER}"
	datePlaceholder        = "{DATE}"
	placeholderGeneric     = "{?}"
	stackFramesPlaceholder = "{StackFrames}"
)

// normalizeMessage processes an error message to extract its stable part using regex.
func normalizeMessage(message string) string {
	if message == "" {
		return ""
	}

	// Check for special patterns that should return empty
	if specialFilePattern.MatchString(strings.TrimSpace(message)) {
		return ""
	}

	normalized := message

	// Apply normalization in priority order

	// 1. URLs (highest priority)
	normalized = urlRegex.ReplaceAllString(normalized, urlPlaceholder)

	// 2. IP addresses with ports
	normalized = ipv4Regex.ReplaceAllString(normalized, ipPlaceholder)
	normalized = ipv6Regex.ReplaceAllString(normalized, ipPlaceholder)

	// 3. JWT tokens (before general placeholders)
	normalized = jwtRegex.ReplaceAllString(normalized, hashPlaceholder)

	// 4. UUIDs
	normalized = uuidRegex.ReplaceAllString(normalized, hashPlaceholder)

	// 5. API keys and tokens (before other patterns)
	normalized = apiKeyRegex.ReplaceAllString(normalized, placeholderGeneric)

	// 6. Base64 encoded strings (before ObjectIDs and general patterns)
	normalized = base64Regex.ReplaceAllString(normalized, hashPlaceholder)

	// 7. MongoDB ObjectIDs (but handle path context specially)
	// First handle ObjectIDs in URL paths as numbers
	objectIDInPathRegex := regexp.MustCompile(`/([0-9a-f]{24})(?:/|$|\|)`)
	normalized = objectIDInPathRegex.ReplaceAllStringFunc(normalized, func(match string) string {
		if strings.HasSuffix(match, "/") {
			return "/" + numberPlaceholder + "/"
		} else if strings.HasSuffix(match, "|") {
			return "/" + numberPlaceholder + "|"
		} else {
			return "/" + numberPlaceholder
		}
	})

	// 8. Date and time patterns (before numbers to avoid breaking dates)
	// Chinese date format first (before numbers)
	normalized = chineseDateRegex.ReplaceAllString(normalized, datePlaceholder)
	// ISO 8601 format
	normalized = iso8601DateRegex.ReplaceAllString(normalized, datePlaceholder)
	// Traditional format
	normalized = traditionalDateRegex.ReplaceAllString(normalized, datePlaceholder)
	// RFC2822 format
	normalized = rfc2822DateRegex.ReplaceAllString(normalized, datePlaceholder)
	// Complete datetime formats (before date only to avoid splitting)
	normalized = completeDateTimeRegex.ReplaceAllString(normalized, datePlaceholder)
	// Month name formats (before date only to prioritize full month names)
	normalized = monthNameRegex.ReplaceAllString(normalized, datePlaceholder)
	// Date only formats
	normalized = dateOnlyRegex.ReplaceAllString(normalized, datePlaceholder)
	// Time only formats (before general numbers)
	normalized = timeOnlyRegex.ReplaceAllString(normalized, datePlaceholder)

	// 9. All numbers (after dates, before hex patterns to prioritize number interpretation)
	// This includes both large and small numbers
	allNumberRegex := regexp.MustCompile(`\b\d+\b`)
	normalized = allNumberRegex.ReplaceAllString(normalized, numberPlaceholder)

	// 10. Trace IDs (32 hex chars, before other hex patterns)
	normalized = traceIDRegex.ReplaceAllString(normalized, hashPlaceholder)

	// 11. Hash values (SHA1: 40, SHA256: 64)
	normalized = hashRegex.ReplaceAllString(normalized, hashPlaceholder)

	// 12. Remaining ObjectIDs as hashes (after numbers, before task IDs)
	normalized = objectIDRegex.ReplaceAllString(normalized, hashPlaceholder)

	// 13. Long hex strings (after specific patterns, before domains)
	normalized = hexLongRegex.ReplaceAllString(normalized, hashPlaceholder)

	// 14. Long identifiers with underscores (before domains)
	normalized = longIdentifierRegex.ReplaceAllString(normalized, hashPlaceholder)

	// 15. Domain names (before task IDs and file paths to avoid conflicts)
	normalized = domainRegex.ReplaceAllString(normalized, domainPlaceholder)

	// 16. Source file paths (after domains to avoid .cc conflicts)
	normalized = sourceFilePathRegex.ReplaceAllString(normalized, filePlaceholder)

	// 17. Task IDs and similar identifiers (after domains and files)
	normalized = taskIdRegex.ReplaceAllString(normalized, placeholderGeneric)

	// 18. Short uppercase identifiers in specific contexts (like key:VALUE, currencyList=VALUE, 平台key：VALUE)
	shortContextRegex := regexp.MustCompile(`((?:key|currency|platform|币种|平台key)(?:Key|List|Name)?[:：=])\s*([A-Z]{3,8})\b`)
	normalized = shortContextRegex.ReplaceAllString(normalized, "${1}"+placeholderGeneric)

	// 19. Payment method identifiers in parentheses (like (PIX), (BRL))
	paymentMethodRegex := regexp.MustCompile(`\(([A-Z]{2,8})\)`)
	normalized = paymentMethodRegex.ReplaceAllString(normalized, "({?})")

	// 20. Handle stack traces
	stackStartIndex := strings.Index(normalized, stackTraceStartPattern)
	if stackStartIndex != -1 {
		// Keep the part before the stack trace starts
		normalized = normalized[:stackStartIndex]
		// Append the placeholder
		normalized += " " + stackFramesPlaceholder
	}

	// Final cleanup
	normalized = multiSpaceRegex.ReplaceAllString(normalized, " ")
	normalized = strings.TrimSpace(normalized)

	return normalized
}
