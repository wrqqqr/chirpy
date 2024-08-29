package main

import "testing"

func TestReplacementBadWords(t *testing.T) {
	var cases = []struct {
		input string
		want  string
	}{
		{"This is a kerfuffle opinion I need to share with the world", "This is a **** opinion I need to share with the world"},
		{"kerfuffle sharbert fornax", "**** **** ****"},
		{"This is a kerfuffle sharbert I need fornax share with the world", "This is a **** **** I need **** share with the world"},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			ans := badWordReplacement(tt.input)
			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		})
	}
}
