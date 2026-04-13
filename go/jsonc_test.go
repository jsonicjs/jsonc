/* Copyright (c) 2021-2025 Richard Rodger, MIT License */

package jsonc

import (
	"reflect"
	"strings"
	"testing"
)

// assert is a test helper that checks deep equality.
func assert(t *testing.T, name string, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s:\n  got:  %#v\n  want: %#v", name, got, want)
	}
}

func assertError(t *testing.T, name string, err error, contains string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: expected error containing %q, got nil", name, contains)
		return
	}
	if !strings.Contains(err.Error(), contains) {
		t.Errorf("%s: expected error containing %q, got: %v", name, contains, err)
	}
}

func TestHappy(t *testing.T) {
	r, err := Parse(`{"a":1}`)
	if err != nil {
		t.Fatal(err)
	}
	assert(t, "basic", r, map[string]any{"a": float64(1)})
}

func TestComments(t *testing.T) {
	r, err := Parse("// this is a comment")
	if err != nil {
		t.Fatal(err)
	}
	assert(t, "single-line", r, nil)

	r, err = Parse("// this is a comment\n")
	if err != nil {
		t.Fatal(err)
	}
	assert(t, "single-line-newline", r, nil)

	r, err = Parse("/* this is a comment*/")
	if err != nil {
		t.Fatal(err)
	}
	assert(t, "block", r, nil)

	r, err = Parse("/* this is a \r\ncomment*/")
	if err != nil {
		t.Fatal(err)
	}
	assert(t, "block-crlf", r, nil)

	r, err = Parse("/* this is a \ncomment*/")
	if err != nil {
		t.Fatal(err)
	}
	assert(t, "block-lf", r, nil)

	_, err = Parse("/* this is a")
	assertError(t, "unterminated-block", err, "unterminated_comment")

	_, err = Parse("/* this is a \ncomment")
	assertError(t, "unterminated-block-multiline", err, "unterminated_comment")

	_, err = Parse("/ ttt")
	assertError(t, "invalid-comment", err, "unexpected")
}

func TestStrings(t *testing.T) {
	r, err := Parse(`"test"`)
	if err != nil {
		t.Fatal(err)
	}
	assert(t, "simple", r, "test")

	r, _ = Parse(`"\""`)
	assert(t, "escape-quote", r, `"`)

	r, _ = Parse(`"\/"`)
	assert(t, "escape-slash", r, "/")

	r, _ = Parse(`"\b"`)
	assert(t, "escape-backspace", r, "\b")

	r, _ = Parse(`"\f"`)
	assert(t, "escape-formfeed", r, "\f")

	r, _ = Parse(`"\n"`)
	assert(t, "escape-newline", r, "\n")

	r, _ = Parse(`"\r"`)
	assert(t, "escape-return", r, "\r")

	r, _ = Parse(`"\t"`)
	assert(t, "escape-tab", r, "\t")

	r, _ = Parse(`"\u00DC"`)
	assert(t, "unicode", r, "\u00DC")

	// Note: \v is accepted by the jsonic Go string matcher as a built-in escape.
	// This is a minor deviation from strict JSONC spec which only allows
	// \", \\, \/, \b, \f, \n, \r, \t, and \uXXXX.

	_, err = Parse(`"test`)
	assertError(t, "unterminated", err, "unterminated_string")
}

func TestNumbers(t *testing.T) {
	r, _ := Parse("0")
	assert(t, "zero", r, float64(0))

	r, _ = Parse("0.1")
	assert(t, "decimal", r, 0.1)

	r, _ = Parse("-0.1")
	assert(t, "neg-decimal", r, -0.1)

	r, _ = Parse("-1")
	assert(t, "neg", r, float64(-1))

	r, _ = Parse("1")
	assert(t, "one", r, float64(1))

	r, _ = Parse("123456789")
	assert(t, "large", r, float64(123456789))

	r, _ = Parse("10")
	assert(t, "ten", r, float64(10))

	r, _ = Parse("90")
	assert(t, "ninety", r, float64(90))

	r, _ = Parse("90E+123")
	assert(t, "sci-upper-plus", r, 90E+123)

	r, _ = Parse("90e+123")
	assert(t, "sci-lower-plus", r, 90e+123)

	r, _ = Parse("90e-123")
	assert(t, "sci-lower-minus", r, 90e-123)

	r, _ = Parse("90E-123")
	assert(t, "sci-upper-minus", r, 90E-123)

	r, _ = Parse("90E123")
	assert(t, "sci-upper", r, 90E123)

	r, _ = Parse("90e123")
	assert(t, "sci-lower", r, 90e123)

	_, err := Parse("-")
	if err == nil {
		t.Error("expected error for bare minus")
	}

	_, err = Parse(".0")
	if err == nil {
		t.Error("expected error for leading dot number")
	}
}

func TestKeywords(t *testing.T) {
	r, _ := Parse("true")
	assert(t, "true", r, true)

	r, _ = Parse("false")
	assert(t, "false", r, false)

	r, _ = Parse("null")
	assert(t, "null", r, nil)

	_, err := Parse("True")
	if err == nil {
		t.Error("expected error for capitalized True")
	}

	r, _ = Parse("false//hello")
	assert(t, "value-with-comment", r, false)
}

func TestTrivia(t *testing.T) {
	r, _ := Parse(" ")
	assert(t, "space", r, nil)

	r, _ = Parse("  \t  ")
	assert(t, "tabs", r, nil)

	r, _ = Parse("  \t  \n  \t  ")
	assert(t, "tabs-newlines", r, nil)

	r, _ = Parse("\r\n")
	assert(t, "crlf", r, nil)

	r, _ = Parse("\r")
	assert(t, "cr", r, nil)

	r, _ = Parse("\n")
	assert(t, "lf", r, nil)

	r, _ = Parse("\n\r")
	assert(t, "lfcr", r, nil)

	r, _ = Parse("\n   \n")
	assert(t, "newlines-spaces", r, nil)
}

func TestLiterals(t *testing.T) {
	r, _ := Parse("true")
	assert(t, "true", r, true)

	r, _ = Parse("false")
	assert(t, "false", r, false)

	r, _ = Parse("null")
	assert(t, "null", r, nil)

	r, _ = Parse(`"foo"`)
	assert(t, "string", r, "foo")

	r, _ = Parse(`"\"-\\-\/-\b-\f-\n-\r-\t"`)
	assert(t, "escapes", r, "\"-\\-/-\b-\f-\n-\r-\t")

	r, _ = Parse(`"\u00DC"`)
	assert(t, "unicode", r, "\u00DC")

	r, _ = Parse("9")
	assert(t, "nine", r, float64(9))

	r, _ = Parse("-9")
	assert(t, "neg-nine", r, float64(-9))

	r, _ = Parse("0.129")
	assert(t, "decimal", r, 0.129)

	r, _ = Parse("23e3")
	assert(t, "sci", r, 23e3)

	r, _ = Parse("1.2E+3")
	assert(t, "sci-plus", r, 1.2E+3)

	r, _ = Parse("1.2E-3")
	assert(t, "sci-minus", r, 1.2E-3)

	r, _ = Parse("1.2E-3 // comment")
	assert(t, "num-comment", r, 1.2E-3)
}

func TestObjects(t *testing.T) {
	r, _ := Parse("{}")
	assert(t, "empty", r, map[string]any{})

	r, _ = Parse(`{ "foo": true }`)
	assert(t, "one-field", r, map[string]any{"foo": true})

	r, _ = Parse(`{ "bar": 8, "xoo": "foo" }`)
	assert(t, "two-fields", r, map[string]any{"bar": float64(8), "xoo": "foo"})

	r, _ = Parse(`{ "hello": [], "world": {} }`)
	assert(t, "empty-nested", r, map[string]any{"hello": []any{}, "world": map[string]any{}})

	r, _ = Parse(`{ "a": false, "b": true, "c": [ 7.4 ] }`)
	assert(t, "mixed", r, map[string]any{"a": false, "b": true, "c": []any{7.4}})

	r, _ = Parse(`{ "hello": { "again": { "inside": 5 }, "world": 1 }}`)
	assert(t, "deep-nested", r, map[string]any{
		"hello": map[string]any{
			"again": map[string]any{"inside": float64(5)},
			"world": float64(1),
		},
	})

	r, _ = Parse(`{ "foo": /*hello*/true }`)
	assert(t, "comment-in-obj", r, map[string]any{"foo": true})

	r, _ = Parse(`{ "": true }`)
	assert(t, "empty-key", r, map[string]any{"": true})
}

func TestArrays(t *testing.T) {
	r, _ := Parse("[]")
	assert(t, "empty", r, []any{})

	r, _ = Parse("[ [],  [ [] ]]")
	assert(t, "nested-empty", r, []any{[]any{}, []any{[]any{}}})

	r, _ = Parse("[ 1, 2, 3 ]")
	assert(t, "numbers", r, []any{float64(1), float64(2), float64(3)})

	r, _ = Parse(`[ { "a": null } ]`)
	assert(t, "obj-in-array", r, []any{map[string]any{"a": nil}})
}

func TestObjectErrors(t *testing.T) {
	_, err := Parse("{,}")
	if err == nil {
		t.Error("expected error for leading comma in object")
	}

	_, err = Parse(`{ "foo": true, }`)
	if err == nil {
		t.Error("expected error for trailing comma in object (default)")
	}

	_, err = Parse(`{ "bar": 8 "xoo": "foo" }`)
	if err == nil {
		t.Error("expected error for missing comma in object")
	}

	_, err = Parse(`{ ,"bar": 8 }`)
	if err == nil {
		t.Error("expected error for leading comma")
	}

	_, err = Parse(`{ "bar": 8, "foo": }`)
	if err == nil {
		t.Error("expected error for missing value")
	}

	_, err = Parse(`{ 8, "foo": 9 }`)
	if err == nil {
		t.Error("expected error for number as key")
	}
}

func TestArrayErrors(t *testing.T) {
	_, err := Parse("[,]")
	if err == nil {
		t.Error("expected error for leading comma in array")
	}

	_, err = Parse("[ 1 2, 3 ]")
	if err == nil {
		t.Error("expected error for missing comma in array")
	}

	_, err = Parse("[ ,1, 2, 3 ]")
	if err == nil {
		t.Error("expected error for leading comma in array")
	}

	_, err = Parse("[ ,1, 2, 3, ]")
	if err == nil {
		t.Error("expected error for commas in array")
	}
}

func TestErrors(t *testing.T) {
	_, err := Parse("1,1")
	if err == nil {
		t.Error("expected error for extra content after value")
	}

	_, err = Parse("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestDisallowComments(t *testing.T) {
	nc := MakeJsonic(JsoncOptions{DisallowComments: boolPtr(true)})

	r, err := nc.Parse(`[ 1, 2, null, "foo" ]`)
	if err != nil {
		t.Fatal(err)
	}
	assert(t, "array", r, []any{float64(1), float64(2), nil, "foo"})

	r, err = nc.Parse(`{ "hello": [], "world": {} }`)
	if err != nil {
		t.Fatal(err)
	}
	assert(t, "object", r, map[string]any{"hello": []any{}, "world": map[string]any{}})

	_, err = nc.Parse(`{ "foo": /*comment*/ true }`)
	if err == nil {
		t.Error("expected error for comment when comments are disallowed")
	}
}

func TestTrailingComma(t *testing.T) {
	jc := MakeJsonic(JsoncOptions{AllowTrailingComma: boolPtr(true)})

	r, err := jc.Parse(`{ "hello": [], }`)
	if err != nil {
		t.Fatal(err)
	}
	assert(t, "obj-trailing", r, map[string]any{"hello": []any{}})

	r, err = jc.Parse(`{ "hello": [] }`)
	if err != nil {
		t.Fatal(err)
	}
	assert(t, "obj-no-trailing", r, map[string]any{"hello": []any{}})

	r, err = jc.Parse(`{ "hello": [], "world": {}, }`)
	if err != nil {
		t.Fatal(err)
	}
	assert(t, "obj-multi-trailing", r, map[string]any{"hello": []any{}, "world": map[string]any{}})

	r, err = jc.Parse(`[ 1, 2, ]`)
	if err != nil {
		t.Fatal(err)
	}
	assert(t, "arr-trailing", r, []any{float64(1), float64(2)})

	r, err = jc.Parse(`[ 1, 2 ]`)
	if err != nil {
		t.Fatal(err)
	}
	assert(t, "arr-no-trailing", r, []any{float64(1), float64(2)})

	// Default parser should reject trailing commas.
	j := MakeJsonic()

	_, err = j.Parse(`{ "hello": [], }`)
	if err == nil {
		t.Error("expected error for trailing comma with default options")
	}

	_, err = j.Parse(`[ 1, 2, ]`)
	if err == nil {
		t.Error("expected error for trailing comma in array with default options")
	}
}

func TestMisc(t *testing.T) {
	j := MakeJsonic()

	r, _ := j.Parse(`{ "foo": "bar" }`)
	assert(t, "simple-obj", r, map[string]any{"foo": "bar"})

	r, _ = j.Parse(`{ "foo": {"bar": 1, "car": 2 } }`)
	assert(t, "nested-obj", r, map[string]any{
		"foo": map[string]any{"bar": float64(1), "car": float64(2)},
	})

	r, _ = j.Parse(`{ "foo": {"bar": 1, "car": 8 }, "goo": {} }`)
	assert(t, "multi-nested", r, map[string]any{
		"foo": map[string]any{"bar": float64(1), "car": float64(8)},
		"goo": map[string]any{},
	})

	_, err := j.Parse(`{ "dep": {"bar": 1, "car": `)
	if err == nil {
		t.Error("expected error for unterminated object")
	}

	_, err = j.Parse(`{ "dep": {"bar": 1,, "car": `)
	if err == nil {
		t.Error("expected error for double comma")
	}

	_, err = j.Parse(`{ "dep": {"bar": "na", "dar": "ma", "car":  } }`)
	if err == nil {
		t.Error("expected error for missing value")
	}

	r, _ = j.Parse(`["foo", null ]`)
	assert(t, "arr-mixed", r, []any{"foo", nil})

	_, err = j.Parse(`["foo", null, ]`)
	if err == nil {
		t.Error("expected error for trailing comma in array")
	}

	_, err = j.Parse(`["foo", null,, ]`)
	if err == nil {
		t.Error("expected error for double comma in array")
	}

	r, _ = j.Parse("true")
	assert(t, "bare-true", r, true)

	r, _ = j.Parse("false")
	assert(t, "bare-false", r, false)

	r, _ = j.Parse("null")
	assert(t, "bare-null", r, nil)

	r, _ = j.Parse("23")
	assert(t, "bare-num", r, float64(23))

	r, _ = j.Parse("-1.93e-19")
	assert(t, "sci-notation", r, -1.93e-19)

	r, _ = j.Parse(`"hello"`)
	assert(t, "bare-string", r, "hello")

	r, _ = j.Parse("[]")
	assert(t, "empty-arr", r, []any{})

	r, _ = j.Parse("[ 1 ]")
	assert(t, "single-arr", r, []any{float64(1)})

	r, _ = j.Parse(`[ 1, "x"]`)
	assert(t, "mixed-arr", r, []any{float64(1), "x"})

	r, _ = j.Parse("[[]]")
	assert(t, "nested-arr", r, []any{[]any{}})

	r, _ = j.Parse("{ }")
	assert(t, "empty-obj", r, map[string]any{})

	r, _ = j.Parse(`{ "val": 1 }`)
	assert(t, "val-obj", r, map[string]any{"val": float64(1)})

	r, _ = j.Parse(`{"id": "$", "v": [ null, null] }`)
	assert(t, "complex-obj", r, map[string]any{"id": "$", "v": []any{nil, nil}})

	_, err = j.Parse(`{  "id": { "foo": { } } , }`)
	if err == nil {
		t.Error("expected error for trailing comma")
	}

	r, _ = j.Parse(`{ "foo": { "goo": 3 } }`)
	assert(t, "nested-num", r, map[string]any{"foo": map[string]any{"goo": float64(3)}})

	r, _ = j.Parse("[\r\n0,\r\n1,\r\n2\r\n]")
	assert(t, "crlf-arr", r, []any{float64(0), float64(1), float64(2)})

	r, _ = j.Parse(`/* g */ { "foo": //f` + "\n" + `"bar" }`)
	assert(t, "comments-mixed", r, map[string]any{"foo": "bar"})

	r, _ = j.Parse("/* g\r\n */ { \"foo\": //f\n\"bar\" }")
	assert(t, "comments-crlf", r, map[string]any{"foo": "bar"})

	r, _ = j.Parse("/* g\n */ { \"foo\": //f\n\"bar\"\n}")
	assert(t, "comments-lf", r, map[string]any{"foo": "bar"})

	r, _ = j.Parse(`{ "key1": { "key11": [ "val111", "val112" ] }, "key2": [ { "key21": false, "key22": 221 }, null, [{}] ] }`)
	assert(t, "complex", r, map[string]any{
		"key1": map[string]any{"key11": []any{"val111", "val112"}},
		"key2": []any{
			map[string]any{"key21": false, "key22": float64(221)},
			nil,
			[]any{map[string]any{}},
		},
	})
}

func TestUsePlugin(t *testing.T) {
	j := MakeJsonic()
	result, err := j.Parse(`{"a": 1, "b": "hello"}`)
	if err != nil {
		t.Fatal(err)
	}
	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", result)
	}
	assert(t, "plugin", m, map[string]any{"a": float64(1), "b": "hello"})
}
