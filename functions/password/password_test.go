package password

import "testing"

func TestPass(t *testing.T) {
	h := NewPassword("test password")
	if !CheckPassword("test password", h) {
		t.Errorf("password doesn't match")
	}

	if CheckPassword("different password", h) {
		t.Errorf("different password matches")
	}
}
