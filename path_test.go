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
	}

	for _, d := range data {
		res, err := findShortestPath(d.input)
		if err != nil {
			t.Errorf("Expected no but found error;  %v \n", err)
		}

		if !reflect.DeepEqual(res, d.output) {
			t.Errorf("The two arrays; (%v) and (%v) are not equal", res, d.output)
		}
	}
}
