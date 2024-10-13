package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/yavurb/rill/testhelpers"
)

// Mock object for testing
type MockObject struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

// Helper function to create a base64 encoded string
func createBase64EncodedString(obj MockObject, compress bool) string {
	b, _ := json.Marshal(obj)
	if compress {
		b = zip(b)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func TestEncode(t *testing.T) {
	t.Run("it should encode a struct object to a json base64 string", func(t *testing.T) {
		obj := MockObject{Field1: "John Doe", Field2: 30}

		encoded, err := Encode(obj, false)
		if err != nil {
			t.Errorf("Error encoding object: %v", err)
		}

		// Marshal the object to JSON
		want, err := json.Marshal(obj)
		if err != nil {
			t.Errorf("Failed to marshal object: %v", err)
		}

		// Encode the bytes to base64
		got, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			t.Errorf("Failed to decode base64. Got error: %v", err)
		}

		if !testhelpers.CompareByteMaps(want, got) {
			t.Errorf("Mismatch encoding struct:\n%s", cmp.Diff(want, got))
		}
	})
}

func TestDecode(t *testing.T) {
	t.Run("Normal decoding without compression", func(t *testing.T) {
		obj := MockObject{Field1: "test", Field2: 123}
		encodedStr := createBase64EncodedString(obj, false)
		var decodedObj MockObject

		Decode(encodedStr, &decodedObj, false)

		if decodedObj != obj {
			t.Errorf("Expected %v, got %v", obj, decodedObj)
		}
	})

	t.Run("Decoding with compression", func(t *testing.T) {
		compress := true
		obj := MockObject{Field1: "test", Field2: 123}
		encodedStr := createBase64EncodedString(obj, compress)
		var decodedObj MockObject

		fmt.Println(encodedStr)

		err := Decode(encodedStr, &decodedObj, compress)
		if err != nil {
			t.Errorf("Error decoding object: %v", err)
		}

		if decodedObj != obj {
			t.Errorf("Expected %v, got %v", obj, decodedObj)
		}
	})

	t.Run("Invalid base64 string", func(t *testing.T) {
		var decodedObj MockObject
		err := Decode("invalid_base64", &decodedObj, false)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}
