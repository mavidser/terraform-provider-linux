package linux

import (
	"testing"
)

func TestInvalidPath(t *testing.T) {
	_, err := validatePath("some/path", "")
	if err == nil {
		t.Errorf("Non-absolute path should be invalid: %v", err)
	}
}

func TestValidPath(t *testing.T) {
	_, err := validatePath("/some/path", "")
	if err != nil {
		t.Errorf("Absolute path should valid: %v", err)
	}
}

func TestValidOwner(t *testing.T) {
	_, err := validateOwner("owner:group", "")
	if err != nil {
		t.Errorf("owner:group is a valid definition: %v", err)
	}

	_, err = validateOwner("123:123", "")
	if err != nil {
		t.Errorf("123:123 is a valid definition: %v", err)
	}
}

func TestInvalidOwner(t *testing.T) {
	_, err := validateOwner("owner", "")
	if err == nil {
		t.Errorf("Owners should be formatted in user:group form: %v", err)
	}
	_, err = validateOwner(123, "")
	if err == nil {
		t.Errorf("Owners should be formatted in user:group form: %v", err)
	}
	_, err = validateOwner(":group", "")
	if err == nil {
		t.Errorf("Owners should be formatted in user:group form: %v", err)
	}
	_, err = validateOwner("user:", "")
	if err == nil {
		t.Errorf("Owners should be formatted in user:group form: %v", err)
	}
}
