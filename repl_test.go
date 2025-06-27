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
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "   ONE    Two    ThRee  ",
			expected: []string{"one", "two", "three"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("For input %q, expected length %d but got %d", c.input, len(c.expected), len(actual))
			continue
		}

		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("For input %q, at index %d expected %q but got %q", c.input, i, c.expected[i], actual[i])
			}
		}
	}
}
