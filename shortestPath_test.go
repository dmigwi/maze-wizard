package main

import (
	"reflect"
	"testing"
)

type testdata struct {
	input  []byte
	output []string
}

func TestProcessShortestPath(t *testing.T) {
	data := []testdata{
		{
			input:  []byte(`{"forward": "tiger", "left": {"forward": {"upstairs": "exit"}, "left": "dragon"}, "right": {"forward":"dead end"}}`),
			output: []string{"left", "forward", "upstairs"},
		},
		{
			input:  []byte(`{"forward": "exit"}`),
			output: []string{"forward"},
		},
		{
			input:  []byte(`{"forward": "tiger", "left": "ogre", "right": "demon"}`),
			output: []string{},
		},
		{
			input:  []byte(`{"left": "tiger", "forward": {"forward": "exit"}, "right": "exit"}`),
			output: []string{"right"},
		},
	}

	for _, d := range data {
		res, err := findShortestPath(d.input)
		if err != nil {
			t.Fatalf("Expected no but found error;  %v \n", err)
		}

		if !reflect.DeepEqual(res, d.output) {
			t.Fatalf("The two arrays; (%v) and (%v) are not equal", res, d.output)
		}
	}
}
