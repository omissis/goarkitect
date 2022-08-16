package santhosh_test

import (
	"goarkitect/internal/schema/santhosh"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/santhosh-tekuri/jsonschema"
)

func Test_JoinPtrPath(t *testing.T) {
	testCases := []struct {
		desc string
		path []any
		want string
	}{
		{
			desc: "empty path",
			path: []any{},
			want: "#",
		},
		{
			desc: "strings-only path",
			path: []any{"foo", "bar"},
			want: "#/foo/bar",
		},
		{
			desc: "ints-only path",
			path: []any{0, 1, 2, 3},
			want: "#/0/1/2/3",
		},
		{
			desc: "strings and ints path",
			path: []any{"foo", 3, "bar", 5},
			want: "#/foo/3/bar/5",
		},
		{
			desc: "no-strings and no-ints path",
			path: []any{0.0, 'a', true},
			want: "#",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := santhosh.JoinPtrPath(tC.path)
			if got != tC.want {
				t.Errorf("got %v, want %v", got, tC.want)
			}
		})
	}
}

func Test_GetValueAtPath(t *testing.T) {
	testCases := []struct {
		desc    string
		obj     any
		path    []any
		want    any
		wantErr error
	}{
		{
			desc: "1-level obj",
			obj:  map[string]any{"foo": "bar"},
			path: []any{"foo"},
			want: "bar",
		},
		{
			desc: "2-levels obj",
			obj: map[string]any{
				"foo": map[string]any{
					"bar": 123,
				},
			},
			path: []any{"foo", "bar"},
			want: 123,
		},
		{
			desc: "2-levels obj, mixed data types",
			obj: []any{
				map[string]any{"foo": []any{"foo", 123}},
			},
			path: []any{0, "foo", 1},
			want: 123,
		},
		{
			desc: "2-levels obj, mixed data types",
			obj: map[string]any{
				"foo": []any{
					map[string]any{"bar": 123},
				},
			},
			path: []any{"foo", 0, "bar"},
			want: 123,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := santhosh.GetValueAtPath(tC.obj, tC.path)
			if got != tC.want {
				t.Errorf("got %v, want %v", got, tC.want)
			}
			if err != tC.wantErr {
				t.Errorf("got error %v, want %v", err, tC.wantErr)
			}
		})
	}
}

func Test_GetPtrPaths(t *testing.T) {
	testCases := []struct {
		desc string
		err  *jsonschema.ValidationError
		want [][]any
	}{
		{
			desc: "get the value pointed by instance ptr",
			err: &jsonschema.ValidationError{
				Message:     "missing properties: \"filePath\"",
				InstancePtr: "#/rules/0/matcher",
				SchemaURL:   "/home/user/api/config_schema.json",
				SchemaPtr:   "#/definitions/fileMatcherOne/required",
			},
			want: [][]any{
				{"rules", 0, "matcher"},
			},
		},
		{
			desc: "get the value pointed by instance ptr and its cause",
			err: &jsonschema.ValidationError{
				Message:     "missing properties: \"filePath\"",
				InstancePtr: "#/rules/0/matcher",
				SchemaURL:   "/home/user/api/config_schema.json",
				SchemaPtr:   "#/definitions/fileMatcherOne/required",
				Causes: []*jsonschema.ValidationError{
					{
						Message:     "wrong property: \"filePaths\"",
						InstancePtr: "#/rules/0/matcher/filePaths",
						SchemaURL:   "/home/user/api/config_schema.json",
						SchemaPtr:   "#/definitions/fileMatcherOne/required",
					},
				},
			},
			want: [][]any{
				{"rules", 0, "matcher"},
				{"rules", 0, "matcher", "filePaths"},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := santhosh.GetPtrPaths(tC.err)
			if !cmp.Equal(got, tC.want, cmpopts.EquateEmpty()) {
				t.Errorf("expected %v, got %v", tC.want, got)
			}
		})
	}
}
