package main

import (
	"reflect"
	"testing"
)

var identifierTests = []struct {
	values map[string]map[string]string
	id     Identifier
	found  bool
	value  string
}{
	{
		values: nil,
		id:     Identifier{key: "foo"},
		found:  false,
	},
	{
		values: map[string]map[string]string{"": {"foo": "example"}},
		id:     Identifier{key: "foo"},
		found:  true,
		value:  "example",
	},
	{
		values: map[string]map[string]string{"": {"foo": "example"}},
		id:     Identifier{key: "bar"},
		found:  false,
	},
	{
		values: map[string]map[string]string{"example": {"foo": "example"}},
		id:     Identifier{scope: "example", key: "foo"},
		found:  true,
		value:  "example",
	},
	{
		values: map[string]map[string]string{"": {"foo": "example"}},
		id:     Identifier{scope: "example", key: "foo"},
		found:  false,
	},
	{
		values: map[string]map[string]string{"example": {"foo": "example"}},
		id:     Identifier{scope: "example", key: "bar"},
		found:  false,
	},
}

func Test_found(t *testing.T) {
	for _, test := range identifierTests {
		got := found(test.values, &test.id)
		if got != test.found {
			t.Errorf("found not correct for %+v (found: %+v, got: %+v)", test.id, test.found, got)
		}
	}
}

func Test_collect(t *testing.T) {
	ids := []*Identifier{
		&Identifier{scope: "foo", key: "foo"},
		&Identifier{scope: "foo", key: "bar"},
		&Identifier{scope: "zoo", key: "foo"},
		&Identifier{scope: "foo", key: "foo"},
		&Identifier{scope: "foo", key: "baz"},
		&Identifier{scope: "qux", key: "bar"},
	}
	expectedFoo := &IdentifierGroup{
		scope: "foo",
		keys:  []string{"foo", "bar", "baz"},
	}
	expectedBar := &IdentifierGroup{
		scope: "bar",
		keys:  nil,
	}
	idgFoo := collect(ids, "foo")
	if !reflect.DeepEqual(idgFoo, expectedFoo) {
		t.Errorf("collect not correct (expected: %+v, got: %+v)", expectedFoo, idgFoo)
	}
	idgBar := collect(ids, "bar")
	if !reflect.DeepEqual(idgBar, expectedBar) {
		t.Errorf("collect not correct (expected: %+v, got: %+v)", expectedBar, idgBar)
	}
}

func Test_insert(t *testing.T) {
	values := make(map[string]map[string]string)
	id := &Identifier{key: "foo"}
	value := "bar"
	insert(values, id, value)
	v, ok := values[""]["foo"]
	if !ok {
		t.Errorf("insert failed for %+v", id)
	}
	if v != value {
		t.Errorf("insert not correctly for %+v (found: %+v, got: %+v)", id, v, value)
	}
	id = &Identifier{scope: "foo", key: "bar"}
	value = "example"
	insert(values, id, value)
	v, ok = values["foo"]["bar"]
	if !ok {
		t.Errorf("insert failed for %+v", id)
	}
	if v != value {
		t.Errorf("insert not correctly for %+v (found: %+v, got: %+v)", id, v, value)
	}
}

func Test_empty(t *testing.T) {
	tests := []struct {
		values   map[string]map[string]string
		expected bool
	}{
		{
			values:   nil,
			expected: true,
		},
		{
			values: map[string]map[string]string{
				"foo": map[string]string{},
			},
			expected: true,
		},
		{
			values: map[string]map[string]string{
				"foo": map[string]string{
					"bar": "",
					"baz": "",
				},
			},
			expected: true,
		},
		{
			values: map[string]map[string]string{
				"foo": map[string]string{
					"bar": "",
					"baz": "",
				},
				"bar": map[string]string{
					"qux": "quux",
				},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		got := empty(test.values)
		if got != test.expected {
			t.Errorf("empty not correctly for %+v (expected: %+v, got: %+v)", test.values, test.expected, got)
		}
	}
}

func Test_lookup(t *testing.T) {
	for _, test := range identifierTests {
		got := lookup(test.values, &test.id)
		if got != test.value {
			t.Errorf("lookup not correct for %+v (expected: %+v, got: %+v)", test.id, test.value, got)
		}
	}
}
