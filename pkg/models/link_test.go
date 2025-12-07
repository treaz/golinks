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
			name: "Invalid Keyword - Reserved Word",
			input: LinkInput{
				Keyword:     "api",
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
			name: "Invalid Destination - No Scheme",
			input: LinkInput{
				Keyword:     "google",
				Destination: "google.com",
			},
			wantErr: true,
		},
		{
			name: "Valid Parameterized Destination",
			input: LinkInput{
				Keyword:     "jira",
				Destination: "https://jira.com/browse/{*}",
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
				Keyword:         "google",
				Destination:     "https://google.com",
				Description:     "Search",
				IsParameterized: false,
			},
		},
		{
			name: "Parameterized Conversion - Implicit",
			input: LinkInput{
				Keyword:     "Jira",
				Destination: "https://jira.com/{*}",
			},
			want: Link{
				Keyword:         "jira",
				Destination:     "https://jira.com/{*}",
				IsParameterized: true,
			},
		},
		{
			name: "Parameterized Conversion - Explicit",
			input: LinkInput{
				Keyword:         "Wiki",
				Destination:     "https://wiki.com",
				IsParameterized: true,
			},
			want: Link{
				Keyword:         "wiki",
				Destination:     "https://wiki.com",
				IsParameterized: true,
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
			if got.IsParameterized != tt.want.IsParameterized {
				t.Errorf("ToNative().IsParameterized = %v, want %v", got.IsParameterized, tt.want.IsParameterized)
			}
		})
	}
}
