package i18n_test

import (
	"reflect"
	"testing"
	
	"github.com/N1xx1/golang-i18n"
)

func assertEqual(t *testing.T, expected interface{}, actual interface{}) bool {
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected ", expected, ", got ", actual)
		return false
	}
	return true
}

func assertNil(t *testing.T, actual interface{}) bool {
	if actual != nil {
		t.Error("Expected nil, got ", actual)
		return false
	}
	return true
}

func TestExampleFile(t *testing.T) {
	T, err := i18n.Tfunc("./example_configuration.i18n")
	assertNil(t, err)
	
	assertEqual(t, "Hello World!", T("helloWorld"))
	assertEqual(t, "There are 42 people", T("exampleCount", 42))
}

func TestGlobalT(t *testing.T) {
	err := i18n.GlobalTfunc("./example_configuration.i18n")
	assertNil(t, err)
	
	assertEqual(t, "Hello World!", i18n.T("helloWorld"))
	assertEqual(t, "There are 42 people", i18n.T("exampleCount", 42))
}