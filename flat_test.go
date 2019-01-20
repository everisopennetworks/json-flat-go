package flat

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		given   string
		want    map[string]interface{}
		options Options
	}{
		// test with different primitives
		// String: 'world',
		// Number: 1234.99,
		// Boolean: true,
		// null: null,
		{
			`{"hello": "world"}`,
			map[string]interface{}{"hello": "world"},
			Options{},
		},
		{
			`{"hello": 1234.99}`,
			map[string]interface{}{"hello": 1234.99},
			Options{},
		},
		{
			`{"hello": true}`,
			map[string]interface{}{"hello": true},
			Options{},
		},
		{
			`{"hello": null}`,
			map[string]interface{}{"hello": nil},
			Options{},
		},
		// nested once
		{
			`{"hello":{}}`,
			map[string]interface{}{"hello": map[string]interface{}{}},
			Options{},
		},
		{
			`{"hello":{"world":"good morning"}}`,
			map[string]interface{}{"hello.world": "good morning"},
			Options{},
		},
		{
			`{"hello":{"world":1234.99}}`,
			map[string]interface{}{"hello.world": 1234.99},
			Options{},
		},
		{
			`{"hello":{"world":true}}`,
			map[string]interface{}{"hello.world": true},
			Options{},
		},
		{
			`{"hello":{"world":null}}`,
			map[string]interface{}{"hello.world": nil},
			Options{},
		},
		// empty slice
		{
			`{"hello":{"world":[]}}`,
			map[string]interface{}{"hello.world": []interface{}{}},
			Options{},
		},
		// slice
		{
			`{"hello":{"world":["one","two"]}}`,
			map[string]interface{}{
				"hello.world.0": "one",
				"hello.world.1": "two",
			},
			Options{},
		},
		// nested twice
		{
			`{"hello":{"world":{"again":"good morning"}}}`,
			map[string]interface{}{"hello.world.again": "good morning"},
			Options{},
		},
		// multiple keys
		{
			`{
				"hello": {
					"lorem": {
						"ipsum":"again",
						"dolor":"sit"
					}
				},
				"world": {
					"lorem": {
						"ipsum":"again",
						"dolor":"sit"
					}
				}
			}`,
			map[string]interface{}{
				"hello.lorem.ipsum": "again",
				"hello.lorem.dolor": "sit",
				"world.lorem.ipsum": "again",
				"world.lorem.dolor": "sit"},
			Options{},
		},
		// empty object
		{
			`{"hello":{"empty":{"nested":{}}}}`,
			map[string]interface{}{"hello.empty.nested": map[string]interface{}{}},
			Options{},
		},
		// custom delimiter
		{
			`{"hello":{"world":{"again":"good morning"}}}`,
			map[string]interface{}{"hello:world:again": "good morning"},
			Options{Delimiter: ":"},
		},
		// custom depth
		{
			`{
				"hello": {
					"world": {
						"again": "good morning"
					}
				},
				"lorem": {
					"ipsum": {
						"dolor": "good evening"
					}
				}
			}
			`,
			map[string]interface{}{
				"hello.world": map[string]interface{}{"again": "good morning"},
				"lorem.ipsum": map[string]interface{}{"dolor": "good evening"},
			},
			Options{MaxDepth: 2},
		},
		// custom safe = true
		{
			`{"hello":{"world":["one","two"]}}`,
			map[string]interface{}{
				"hello.world": []interface{}{"one", "two"},
			},
			Options{Safe: true},
		},
	}
	for i, test := range tests {
		var given interface{}
		err := json.Unmarshal([]byte(test.given), &given)
		if err != nil {
			t.Errorf("%d: failed to unmarshal test: %v", i+1, err)
		}
		got, err := Flatten(given.(map[string]interface{}), test.options)
		if err != nil {
			t.Errorf("%d: failed to flatten: %v", i+1, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d: mismatch, got: %v want: %v", i+1, got, test.want)
		}
	}
}

func TestFPrimitive(t *testing.T) {
	tests := []struct {
		prefix  string
		nested  interface{}
		options Options
		want    map[string]interface{}
	}{
		// test with different primitives
		// String: 'good morning',
		// Number: 1234.99,
		// Boolean: true,
		// null: null,
		{
			"hello",
			"world",
			Options{
				MaxDepth: 20,
			},
			map[string]interface{}{"hello": "world"},
		},
		{
			"hello",
			1234.49,
			Options{
				MaxDepth: 20,
			},
			map[string]interface{}{"hello": 1234.49},
		},
		{
			"hello",
			true,
			Options{
				MaxDepth: 20,
			},
			map[string]interface{}{"hello": true},
		},
		{
			"hello",
			nil,
			Options{
				MaxDepth: 20,
			},
			map[string]interface{}{"hello": nil},
		},
	}
	for i, test := range tests {
		got, err := f(test.prefix, 0, test.nested, test.options)
		if err != nil {
			t.Errorf("%d: failed to flatten: %v", i+1, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d: mismatch, got: %v want: %v", i+1, got, test.want)
		}
	}
}
func TestFMap(t *testing.T) {
	tests := []struct {
		prefix  string
		nested  string
		options Options
		want    map[string]interface{}
	}{
		// empty map
		{
			"hello",
			`{}`,
			Options{
				Delimiter: ".",
				MaxDepth:  20,
			},
			map[string]interface{}{"hello": map[string]interface{}{}},
		},
		{
			"hello",
			`{"world": "good morning"}`,
			Options{
				Delimiter: ".",
				MaxDepth:  20,
			},
			map[string]interface{}{"hello.world": "good morning"},
		},
		{
			"",
			`{"world": "good morning"}`,
			Options{
				Delimiter: ".",
				MaxDepth:  20,
			},
			map[string]interface{}{"world": "good morning"},
		},
		// nested twice
		{
			"hello",
			`{"world":{"again":"good morning"}}`,
			Options{
				Delimiter: ".",
				MaxDepth:  20,
			},
			map[string]interface{}{"hello.world.again": "good morning"},
		},
		// multiple key
		{
			"hello",
			`{
				"world": {
					"again": "good morning"
				},
				"ipsum": {
					"dolor": "good evening"
				}
			}`,
			Options{
				Delimiter: ".",
				MaxDepth:  20,
			},
			map[string]interface{}{
				"hello.world.again": "good morning",
				"hello.ipsum.dolor": "good evening",
			},
		},
		// empty slice
		{
			"hello",
			`[]`,
			Options{
				Delimiter: ".",
				Safe:      true,
			},
			map[string]interface{}{"hello": []interface{}{}},
		},
		// slice
		{
			"hello",
			`{"world":["one","two"]}`,
			Options{
				Delimiter: ".",
				MaxDepth:  20,
			},
			map[string]interface{}{
				"hello.world.0": "one",
				"hello.world.1": "two",
			},
		},
		// custom delimiter
		{
			"hello",
			`{"world":{"again":"good morning"}}`,
			Options{
				Delimiter: ":",
				MaxDepth:  20,
			},
			map[string]interface{}{"hello:world:again": "good morning"},
		},
		// custom depth
		{
			"",
			`{
				"hello": {
					"world": {
						"again": "good morning"
					}
				},
				"lorem": {
					"ipsum": {
						"dolor": "good evening"
					}
				}
			}
			`,
			Options{
				Delimiter: ".",
				MaxDepth:  2,
			},
			map[string]interface{}{
				"hello.world": map[string]interface{}{"again": "good morning"},
				"lorem.ipsum": map[string]interface{}{"dolor": "good evening"},
			},
		},
		// safe
		{
			"hello",
			`{"world":["one","two"]}`,
			Options{
				Delimiter: ".",
				MaxDepth:  20,
				Safe:      true,
			},
			map[string]interface{}{
				"hello.world": []interface{}{"one", "two"},
			},
		},
	}
	for i, test := range tests {
		var nested interface{}
		err := json.Unmarshal([]byte(test.nested), &nested)
		if err != nil {
			t.Errorf("%d: failed to unmarshal test: %v", i+1, err)
		}
		got, err := f(test.prefix, 0, nested, test.options)
		if err != nil {
			t.Errorf("%d: failed to flatten: %v", i+1, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d: mismatch, got: %v want: %v", i+1, got, test.want)
		}
	}
}

func TestUnflatten(t *testing.T) {
	tests := []struct {
		flat    map[string]interface{}
		options *Options
		want    map[string]interface{}
	}{
		{
			map[string]interface{}{"hello": "world"},
			nil,
			map[string]interface{}{"hello": "world"},
		},
		{
			map[string]interface{}{"hello": 1234.56},
			nil,
			map[string]interface{}{"hello": 1234.56},
		},
		{
			map[string]interface{}{"hello": true},
			nil,
			map[string]interface{}{"hello": true},
		},
		// nested twice
		{
			map[string]interface{}{"hello.world.again": "good morning"},
			nil,
			map[string]interface{}{
				"hello": map[string]interface{}{
					"world": map[string]interface{}{
						"again": "good morning",
					},
				},
			},
		},
		// multiple keys
		{
			map[string]interface{}{
				"hello.lorem.ipsum": "again",
				"hello.lorem.dolor": "sit",
				"world.lorem.ipsum": "again",
				"world.lorem.dolor": "sit",
				"world":             map[string]interface{}{"greet": "hello"},
			},
			nil,
			map[string]interface{}{
				"hello": map[string]interface{}{
					"lorem": map[string]interface{}{
						"ipsum": "again",
						"dolor": "sit",
					},
				},
				"world": map[string]interface{}{
					"greet": "hello",
					"lorem": map[string]interface{}{
						"ipsum": "again",
						"dolor": "sit",
					},
				},
			},
		},
		// nested objects do not clobber each other
		{
			map[string]interface{}{
				"foo.bar": map[string]interface{}{"t": 123},
				"foo":     map[string]interface{}{"k": 456},
			},
			nil,
			map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": map[string]interface{}{
						"t": 123,
					},
					"k": 456,
				},
			},
		},
		// custom delimiter
		{
			map[string]interface{}{
				"hello world again": "good morning",
			},
			&Options{
				Delimiter: " ",
			},
			map[string]interface{}{
				"hello": map[string]interface{}{
					"world": map[string]interface{}{
						"again": "good morning",
					},
				},
			},
		},
		// do not overwrite
		{
			map[string]interface{}{
				"travis":           "true",
				"travis_build_dir": "/home/foo",
			},
			&Options{
				Delimiter: "_",
			},
			map[string]interface{}{
				"travis": "true",
			},
		},
		// todo
		// overwrite true
		// {
		// 	map[string]interface{}{
		// 		"travis":           "true",
		// 		"travis_build_dir": "/home/foo",
		// 	},
		// 	Options{
		// 		Delimiter: "_",
		// 		Overwrite: true,
		// 	},
		// 	map[string]interface{}{
		// 		"travis": map[string]interface{}{
		// 			"build": map[string]interface{}{
		// 				"dir": "/home/foo",
		// 			},
		// 		},
		// 	},
		// },
	}
	for i, test := range tests {
		got, err := Unflatten(test.flat, test.options)
		if err != nil {
			t.Errorf("%d: failed to unflatten: %v", i+1, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d: mismatch, got: %v want: %v", i+1, got, test.want)
		}
	}
}
