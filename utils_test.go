package toscalib

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type I interface{}

type A struct {
	Greeting string
	Message  string
	Pi       float64
}

type B struct {
	Struct    A
	Ptr       *A
	Answer    int
	Map       map[string]string
	StructMap map[string]interface{}
	Slice     []string
}

func create() I {
	// The type C is actually hidden, but reflection allows us to look inside it
	type C struct {
		String string
	}

	return B{
		Struct: A{
			Greeting: "Hello!",
			Message:  "translate this",
			Pi:       3.14,
		},
		Ptr: &A{
			Greeting: "What's up?",
			Message:  "point here",
			Pi:       3.14,
		},
		Map: map[string]string{
			"Test": "translate this as well",
		},
		StructMap: map[string]interface{}{
			"C": C{
				String: "deep",
			},
		},
		Slice: []string{
			"and one more",
		},
		Answer: 42,
	}
}

func TestNilPointerToStruct(t *testing.T) {
	var original *B
	translated := clone(original)
	if ok := reflect.DeepEqual(original, translated); !ok {
		t.Fatal(spew.Sdump(original), "!=", spew.Sdump(translated))
	}
}

func TestNilPointerToInterface(t *testing.T) {
	var original *I
	translated := clone(original)
	if ok := reflect.DeepEqual(original, translated); !ok {
		t.Fatal(spew.Sdump(original), "!=", spew.Sdump(translated))
	}
}

func TestStructWithNoElements(t *testing.T) {
	type E struct{}
	var original E
	translated := clone(original)
	if ok := reflect.DeepEqual(original, translated); !ok {
		t.Fatal(spew.Sdump(original), "!=", spew.Sdump(translated))
	}
}

func TestEmptyStruct(t *testing.T) {
	var original B
	translated := clone(original)
	if ok := reflect.DeepEqual(original, translated); !ok {
		t.Fatal(spew.Sdump(original), "!=", spew.Sdump(translated))
	}
}

func TestCloneStruct(t *testing.T) {
	created := create()
	original := created.(B)
	translated := clone(original)
	if ok := reflect.DeepEqual(original, translated); !ok {
		t.Fatal(spew.Sdump(original), "!=", spew.Sdump(translated))
	}
}

func TestCloneStructWrappedWithInterface(t *testing.T) {
	created := create()
	original := created
	translated := clone(original)
	if ok := reflect.DeepEqual(original, translated); !ok {
		t.Fatal(spew.Sdump(original), "!=", spew.Sdump(translated))
	}
}

func TestClonePointerToStructWrappedWithInterface(t *testing.T) {
	created := create()
	original := &created
	translated := clone(original)
	if ok := reflect.DeepEqual(original, translated); !ok {
		t.Fatal(spew.Sdump(original), "!=", spew.Sdump(translated))
	}
}

func TestCloneStructWithPointerToStructWrappedWithInterface(t *testing.T) {
	created := create()

	type D struct {
		Payload *I
	}
	original := D{
		Payload: &created,
	}
	translated := clone(original)
	if ok := reflect.DeepEqual(original, translated); !ok {
		t.Fatal(spew.Sdump(original), "!=", spew.Sdump(translated))
	}
}

func TestIsAbsLocalPath(t *testing.T) {
	type localStruct struct {
		fileName string
		expected bool
	}
	testFiles := []localStruct{{"tests/example1.yaml", true},
		{"tests/tosca_helloworld.yaml", true},
		{"does_not_exist.yaml", false},
	}

	dir, _ := os.Getwd()
	for _, testFile := range testFiles {
		result := isAbsLocalPath(filepath.Join(dir, testFile.fileName))
		if result != testFile.expected {
			t.Errorf("checking profile %v, expected %v, actual %v", testFile.fileName, testFile.expected, result)
		}
	}

	//test case relative path
	if isAbsLocalPath("tests/tosca_helloworld.yaml") {
		t.Errorf("relative path test case: expected false but got true")
	}
}
