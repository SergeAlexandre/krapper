package misc

import "testing"

func TestFormatString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HTTPRequestStatus", "HTTP request status"},
		{"APIClient", "API client"},
		{"getHTTPServer", "Get HTTP server"},
		{"hello_world", "Hello world"},
		{"my-XML-parser", "My XML parser"},
		{"ParseJSONFile", "Parse JSON file"},
		{"simpleTest", "Simple test"},
		{"JSON", "JSON"},                      // pure acronym
		{"", ""},                              // empty input
		{"___HTTP__Server", "HTTP server"},    // multiple separators
		{"XMLHTTPRequest", "XMLHTTP request"}, // multiple acronyms combined
		{"UserID", "User ID"},                 // mixed case
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Labelize(tt.input)
			if got != tt.expected {
				t.Errorf("FormatString(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
