package santhosh

import (
	"strconv"
	"strings"

	"github.com/santhosh-tekuri/jsonschema"
	"golang.org/x/exp/slices"
)

func JoinPtrPath(path []any) string {
	strpath := "#"

	for _, key := range path {
		switch v := key.(type) {
		case int:
			strpath += "/" + strconv.Itoa(v)
		case string:
			strpath += "/" + v
		}
	}

	return strpath
}

func GetValueAtPath(obj any, path []any) any {
	for _, key := range path {
		switch v := key.(type) {
		case int:
			obj = obj.([]any)[v]
		case string:
			obj = obj.(map[string]any)[v]
		}
	}

	return obj
}

func GetPtrPaths(err error) [][]any {
	ptrs := extractPtrs(err.(*jsonschema.ValidationError))

	mptrs := minimizePtrs(ptrs)

	return explodePtrs(mptrs)
}

func extractPtrs(err *jsonschema.ValidationError) []string {
	var ptrs []string

	for _, cause := range err.Causes {
		if len(cause.Causes) > 0 {
			ptrs = append(ptrs, extractPtrs(cause)...)
		} else {
			ptrs = append(ptrs, err.InstancePtr)
		}
	}

	return ptrs
}

func minimizePtrs(ptrs []string) []string {
	slices.Sort(ptrs)

	return slices.Compact(ptrs)
}

func explodePtrs(ptrs []string) [][]any {
	eptrs := make([][]any, len(ptrs))
	for i, p := range ptrs {
		eptrs[i] = explodePtr(strings.TrimLeft(p, "#/"))
	}

	return eptrs
}

func explodePtr(ptr string) []any {
	parts := strings.Split(ptr, "/")
	ptrParts := make([]any, len(parts))

	for i, part := range parts {
		if numpart, err := strconv.Atoi(part); err == nil {
			ptrParts[i] = numpart
		} else {
			ptrParts[i] = part
		}
	}

	return ptrParts
}
