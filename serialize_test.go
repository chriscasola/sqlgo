package sqlgo

import (
	"testing"
)

func testSerialize(t *testing.T, thingToSerialize interface{}, expectedVal string, description string) {
	stringValue := Serialize(thingToSerialize)

	if stringValue != expectedVal {
		t.Errorf(`The %v did not serialize correctly, expected "%v" but got "%v"`, description, expectedVal, stringValue)
	}
}

func TestSerialize(t *testing.T) {
	testSerialize(t, "test string", "'test string'", "string")
	testSerialize(t, 34, "34", "integer")
	testSerialize(t, 0, "0", "integer")
	testSerialize(t, -3, "-3", "integer")
	testSerialize(t, -3.2, "-3.2", "float")
	testSerialize(t, 3.2, "3.2", "float")
	testSerialize(t, false, "false", "boolean")
	testSerialize(t, true, "true", "boolean")
	testSerialize(t, nil, "NULL", "nil")
}
