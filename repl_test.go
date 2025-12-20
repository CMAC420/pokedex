package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " CMAC is the best",
			expected: []string{"cmac", "is", "the", "best"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("length error: expected %d, got %d", len(c.expected), len(actual))
			continue
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("word missmatch at index %d: expected %q, got %q", i, c.expected[i], actual[i])
				continue
			}
		}
	}
}
