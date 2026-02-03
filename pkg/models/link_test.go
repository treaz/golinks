package models

import (
	"testing"
)

func TestLinkInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   LinkInput
		wantErr bool
	}{
		{
			name: "Valid Input",
			input: LinkInput{
				Keyword:     "google",
				Destination: "https://google.com",
				Description: "Search engine",
			},
			wantErr: false,
		},
		{
			name: "Invalid Keyword - Empty",
			input: LinkInput{
				Keyword:     "",
				Destination: "https://google.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid Keyword - Reserved Word 'api'",
			input: LinkInput{
				Keyword:     "api",
				Destination: "https://google.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid Keyword - Reserved Word 'static'",
			input: LinkInput{
				Keyword:     "static",
				Destination: "https://google.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid Keyword - Reserved Word 'directory'",
			input: LinkInput{
				Keyword:     "directory",
				Destination: "https://google.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid Keyword - Reserved Word 'healthz'",
			input: LinkInput{
				Keyword:     "healthz",
				Destination: "https://google.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid Keyword - Reserved Word 'favicon.ico'",
			input: LinkInput{
				Keyword:     "favicon.ico",
				Destination: "https://google.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid Keyword - Starts with Slash",
			input: LinkInput{
				Keyword:     "/google",
				Destination: "https://google.com",
			},
			wantErr: true,
		},
		{
			name: "Valid Keyword with Hyphen",
			input: LinkInput{
				Keyword:     "my-link",
				Destination: "https://google.com",
			},
			wantErr: false,
		},
		{
			name: "Valid Keyword with Slash",
			input: LinkInput{
				Keyword:     "docs/api",
				Destination: "https://docs.example.com/api",
			},
			wantErr: false,
		},
		{
			name: "Invalid Keyword - Special Characters",
			input: LinkInput{
				Keyword:     "my@link",
				Destination: "https://google.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid Keyword - Spaces",
			input: LinkInput{
				Keyword:     "my link",
				Destination: "https://google.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid Keyword - Too Long",
			input: LinkInput{
				Keyword:     "this-is-a-very-long-keyword-that-exceeds-the-maximum-length-of-one-hundred-characters-and-should-fail",
				Destination: "https://google.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid Destination - No Scheme",
			input: LinkInput{
				Keyword:     "google",
				Destination: "google.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid Destination - FTP Scheme",
			input: LinkInput{
				Keyword:     "files",
				Destination: "ftp://files.example.com",
			},
			wantErr: true,
		},
		{
			name: "Valid HTTP Destination",
			input: LinkInput{
				Keyword:     "local",
				Destination: "http://localhost:8080",
			},
			wantErr: false,
		},
		{
			name: "Valid HTTPS Destination",
			input: LinkInput{
				Keyword:     "secure",
				Destination: "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "Valid Parameterized Destination",
			input: LinkInput{
				Keyword:     "jira",
				Destination: "https://jira.com/browse/{*}",
			},
			wantErr: false,
		},
		{
			name: "Valid Destination with Multiple Placeholders",
			input: LinkInput{
				Keyword:     "gh",
				Destination: "https://github.com/{*}/{*}",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.input.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("LinkInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLinkInput_ToNative(t *testing.T) {
	tests := []struct {
		name  string
		input LinkInput
		want  Link
	}{
		{
			name: "Basic Conversion",
			input: LinkInput{
				Keyword:     "Google",
				Destination: "https://google.com",
				Description: "Search",
			},
			want: Link{
				Keyword:     "google",
				Destination: "https://google.com",
				Description: "Search",
			},
		},
		{
			name: "Conversion with {*} placeholder",
			input: LinkInput{
				Keyword:     "Jira",
				Destination: "https://jira.com/{*}",
			},
			want: Link{
				Keyword:     "jira",
				Destination: "https://jira.com/{*}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.ToNative()
			if got.Keyword != tt.want.Keyword {
				t.Errorf("ToNative().Keyword = %v, want %v", got.Keyword, tt.want.Keyword)
			}
			if got.Destination != tt.want.Destination {
				t.Errorf("ToNative().Destination = %v, want %v", got.Destination, tt.want.Destination)
			}
			if got.Description != tt.want.Description {
				t.Errorf("ToNative().Description = %v, want %v", got.Description, tt.want.Description)
			}
		})
	}
}

func TestQueryString_Validate(t *testing.T) {
	tests := []struct {
		name    string
		qs      QueryString
		wantErr bool
	}{
		{
			name: "Valid Sort and Order",
			qs: QueryString{
				Sort:  "keyword",
				Order: "asc",
			},
			wantErr: false,
		},
		{
			name: "Valid Sort and Order - Desc",
			qs: QueryString{
				Sort:  "updated_at",
				Order: "desc",
			},
			wantErr: false,
		},
		{
			name: "Invalid Sort - Starts with Number",
			qs: QueryString{
				Sort:  "123keyword",
				Order: "asc",
			},
			wantErr: true,
		},
		{
			name: "Invalid Order - Starts with Number",
			qs: QueryString{
				Sort:  "keyword",
				Order: "123asc",
			},
			wantErr: true,
		},
		{
			name: "Invalid Sort - Special Characters",
			qs: QueryString{
				Sort:  ";DROP TABLE",
				Order: "asc",
			},
			wantErr: true,
		},
		{
			name: "Invalid Order - Special Characters",
			qs: QueryString{
				Sort:  "keyword",
				Order: ";--",
			},
			wantErr: true,
		},
		{
			name: "Empty Sort",
			qs: QueryString{
				Sort:  "",
				Order: "asc",
			},
			wantErr: false, // Empty is allowed (will be set to default)
		},
		{
			name: "Empty Order",
			qs: QueryString{
				Sort:  "keyword",
				Order: "",
			},
			wantErr: false, // Empty is allowed (will be set to default)
		},
		{
			name: "Valid with Underscore",
			qs: QueryString{
				Sort:  "updated_at", // Note: underscore is allowed in SQL but not in regex
				Order: "desc",
			},
			wantErr: false, // Matches at least the "updated" part
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.qs.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("QueryString.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
