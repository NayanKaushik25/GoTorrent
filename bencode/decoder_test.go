package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestDecodeInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"i42e", 42},
		{"i0e", 0},
		{"i-42e", -42},
	}
	i := 0
	for i < len(tests) {
		reader := bytes.NewReader([]byte(tests[i].input))
		got, err := Decode(reader)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got != tests[i].expected {
			t.Errorf("expected %v, got %v", tests[i].expected, got)
		}
		i++
	}
}

func TestDecodeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"4:spam", "spam"},
		{"3:abc", "abc"},
		{"10:helloworld", "helloworld"},
	}
	for _, test := range tests {
		reader := bytes.NewReader([]byte(test.input))
		got, err := Decode(reader)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != test.expected {
			t.Errorf("expected %q,got %q", test.expected, got)
		}
	}
}

func TestDecodeList(t *testing.T) {
	reader := bytes.NewReader([]byte("l4:spaml4:eggsi42eee"))
	got, err := Decode(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list, ok := got.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", got)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(list))
	}
	if list[0] != "spam" {
		t.Errorf("expected spam, got %v", list[0])
	}
	nested, ok := list[1].([]any)
	if !ok {
		t.Fatalf("expected nested list, got %T", list[1])
	}

	if nested[0] != "eggs" {
		t.Errorf("expected eggs, got %v", nested[0])
	}

	if nested[1] != int64(42) {
		t.Errorf("expected 42, got %v", nested[1])
	}
}

func TestDecodeDictionary(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]any
	}{
		{
			name:  "empty dictionary",
			input: "de",
			want:  map[string]any{},
		},
		{
			name:  "single entry",
			input: "d3:foo3:bare",
			want:  map[string]any{"foo": "bar"},
		},
		{
			name:  "multiple entries",
			input: "d3:foo3:bar5:helloi42ee",
			want: map[string]any{
				"foo":   "bar",
				"hello": int64(42),
			},
		},
		{
			name:  "list as value",
			input: "d4:listl1:a1:bee",
			want: map[string]any{
				"list": []any{"a", "b"},
			},
		},
		{
			name:  "nested dictionary",
			input: "d3:food3:bar3:bazee",
			want: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tt.input))

			got, err := Decode(reader)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			dict, ok := got.(map[string]any)
			if !ok {
				t.Fatalf("expected map[string]any, got %T", got)
			}

			if !reflect.DeepEqual(dict, tt.want) {
				t.Errorf("Decode(%q)\n got: %#v\nwant: %#v",
					tt.input, dict, tt.want)
			}
		})
	}
}
