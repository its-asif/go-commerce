package utils

import (
	"testing"
)

func TestHashAndMatchPass(t *testing.T) {
	pass := "s3cr3t"
	hashed, err := HashPass(pass)
	if err != nil {
		t.Fatalf("HashPass error: %v", err)
	}

	if err := MatchPass(hashed, pass); err != nil {
		t.Fatalf("MatchPass failed for correct password: %v", err)
	}

	if err := MatchPass(hashed, "wrong"); err == nil {
		t.Fatalf("MatchPass should fail for incorrect password")
	}
}
