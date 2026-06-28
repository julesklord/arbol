package main

import (
	"testing"
)

func TestEscapeXML(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "empty string",
			in:   "",
			want: "",
		},
		{
			name: "no special characters",
			in:   "hello world",
			want: "hello world",
		},
		{
			name: "ampersand",
			in:   "a & b",
			want: "a &amp; b",
		},
		{
			name: "less than",
			in:   "a < b",
			want: "a &lt; b",
		},
		{
			name: "greater than",
			in:   "a > b",
			want: "a &gt; b",
		},
		{
			name: "double quote",
			in:   `say "hello"`,
			want: "say &quot;hello&quot;",
		},
		{
			name: "single quote",
			in:   "it's",
			want: "it&apos;s",
		},
		{
			name: "all special characters",
			in:   `<a & b>"c"'d'`,
			want: "&lt;a &amp; b&gt;&quot;c&quot;&apos;d&apos;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeXML(tt.in)
			if got != tt.want {
				t.Errorf("escapeXML(%q) = %q; want %q", tt.in, got, tt.want)
			}
		})
	}
}
