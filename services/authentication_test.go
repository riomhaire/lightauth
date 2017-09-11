package services

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	expectedHash := "9ca8e785a75cdcbfb6317a24efe21b870f0b3d302a5ba62d438b13cb87ea8e64"
	hash := HashPassword("secret", "secret")
	if len(hash) == 0 {
		t.Fatalf("Hash returned empty string for valid parameters")
	}
	if expectedHash != hash {
		t.Fatalf("Hash returned %s did not match expected %s\n", hash, expectedHash)
	}
}

func TestLoadUserFile(t *testing.T) {
	err := LoadUsers("../users.csv")
	if err != nil {
		t.Fatalf("Load users failed with %s\n", err)
	}

	// Check more than 1 user in file
	if len(users) == 0 {
		t.Fatalf("Expected more then 0 users\n")
	}

}

// System has a test user - user test password secret roles 'none'
func TestAuthenticateValidUser(t *testing.T) {
	// Need to set up globals
	p := 100
	s := "secret"
	SessionPeriod = &p
	SessionSecret = &s
	err := LoadUsers("../users.csv")
	if err != nil {
		t.Fatalf("Load users failed with %s\n", err)
	}
	service := new(AuthenticationService)
	args := AuthenticationArgs{"test", "secret"}
	result := AuthenticationReply{}
	// Call
	err = service.Authenticate(nil, &args, &result)
	if err != nil {
		t.Fatalf("Unexpected error returned %s\n", err)
	}
	if len(result.Token) == 0 {
		t.Fatalf("No Token returned for known user/password\n")
	}
}

// System has a test user - user test password secret roles 'none'
func TestAuthenticateInValidUser(t *testing.T) {
	// Need to set up globals
	p := 100
	s := "secret"
	SessionPeriod = &p
	SessionSecret = &s
	err := LoadUsers("../users.csv")
	if err != nil {
		t.Fatalf("Load users failed with %s\n", err)
	}
	service := new(AuthenticationService)
	args := AuthenticationArgs{"test", "test"} // Bad password
	result := AuthenticationReply{}
	// Call and expect error 401
	err = service.Authenticate(nil, &args, &result)
	if err == nil {
		t.Fatalf("Expected error returned %s\n", result)
	}

	if !strings.HasPrefix(err.Error(), "401 ") {
		t.Fatalf("Expected 401 error but got %s\n", err.Error())
	}

}
