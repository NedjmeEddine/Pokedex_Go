package main

import (
	"testing"
)

func TestClenInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  Charmander Bulbasaur PIKACHU  ",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {

		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("the length of the actual slice %d does not match the expected length %d", len(actual), len(c.expected))
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("the actual word %q does not match the expected word %q", word, expectedWord)
			}

		}
	}

}
