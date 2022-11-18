package assert

import (
	"reflect"
	"testing"
)

func Equal(t *testing.T, want, got any) bool {
	if want == nil && got != nil || !reflect.DeepEqual(want, got) {
		t.Errorf("[%s] wanted=%v, got=%v\n", t.Name(), want, got)
		return false
	}
	return true
}

func EqualValues(t *testing.T, want, got any, args ...any) bool {
	if want == nil && got != nil || !reflect.DeepEqual(want, got) {
		t.Errorf("[%s] wanted=%v, got=%v\n", t.Name(), want, got)
		return false
	}
	return true
}
