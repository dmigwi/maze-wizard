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
			input: []byte(`{"forward": "tiger", "left": {"forward": {"upstairs": "exit"}, "left": "dragon"}, "right": {"forward":
		"dead end"}}`),
			output: []string{"left", "forward", "upstairs"},
		},
		{
			input:  []byte(`{“forward”: “exit”}`),
			output: []string{"forward"},
		},
		{
			input:  []byte(`{"forward": "tiger", "left": "ogre", "right": "demon"}`),
			output: []string{},
		},
	}

	for _, data := range data {
		res := processShortestPath(data.input, [][]string{})
		if !reflect.DeepEqual(res, data.output) {
			t.Errorf("The two arrays; (%v) and (%v) are not equal", res, data.output)
		}
	}
}
