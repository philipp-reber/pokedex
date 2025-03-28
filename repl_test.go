package main

import (
	"reflect"
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
			input:    "Pikachu VegetA ",
			expected: []string{"pikachu", "vegeta"}, // Fix expected value
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		// Compare slices properly
		if !reflect.DeepEqual(actual, c.expected) {
			t.Errorf("For input %q, expected %v but got %v", c.input, c.expected, actual)
		}
	}
}
