package misc

import (
	"fmt"
	"regexp"
	"strings"
)

func Labelize(input string) string {
	if input == "" {
		return ""
	}

	// Split before CamelCase and before acronyms followed by lowercase
	// e.g. HTTPServer -> HTTP Server, MyXMLParser -> My XML Parser
	reCamel := regexp.MustCompile(`([A-Z]+)([A-Z][a-z])|([a-z0-9])([A-Z])`)
	input = reCamel.ReplaceAllString(input, `${1}${3} ${2}${4}`)

	// Replace - and _ with spaces
	reDelims := regexp.MustCompile(`[_-]+`)
	input = reDelims.ReplaceAllString(input, " ")

	// Normalize spaces
	words := strings.Fields(input)
	if len(words) == 0 {
		return ""
	}

	// Process words: preserve acronyms, lowercase others
	for i := range words {
		if isAcronym(words[i]) {
			words[i] = strings.ToUpper(words[i])
		} else {
			words[i] = strings.ToLower(words[i])
		}
	}

	// Capitalize ONLY the first letter of the first non-acronym word
	if !isAcronym(words[0]) {
		words[0] = strings.Title(words[0])
	}

	return strings.Join(words, " ")
}

// Detect acronyms (2+ uppercase letters)
func isAcronym(s string) bool {
	return len(s) > 1 && s == strings.ToUpper(s)
}

func main() {
	tests := []string{
		"HTTPRequestStatus",
		"APIClient",
		"getHTTPServer",
		"hello_world",
		"my-XML-parser",
		"ParseJSONFile",
	}

	for _, t := range tests {
		fmt.Printf("%s -> %s\n", t, Labelize(t))
	}
}
