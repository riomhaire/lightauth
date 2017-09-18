package services

import (
	"strings"
	"testing"
)

func TestSessionHappyPathLifeCycle(t *testing.T) {
	// Need to set up globals
	p := 1000
	s := "secret"
	SessionPeriod = &p
	SessionSecret = &s

	service := new(SessionService)
	// Create
	user := User{"test", "test", true, []string{}}
	hash, err := service.NewSession(user)
	if err != nil {
		t.Fatalf("Expected no error - but got %s\n", err.Error())
	}
	if len(hash) == 0 {
		t.Fatalf("Expected hash but got nothing\n")
	}
	// Validate
	args := SessionValidateArgs(hash)
	result := SessionValidateReply(false)
	err = service.Validate(nil, &args, &result)
	if err != nil {
		t.Fatalf("Expected no error - but got %s\n", err.Error())
	}
	// Session should be valid
	if !result {
		t.Fatalf("Expected session valid - but it isnt\n")
	}

	// Details - should fail admin/api required
	detailsArgs := SessionDetailsArgs{hash, hash}
	detailsResult := Session{}
	err = service.Details(nil, &detailsArgs, &detailsResult)
	if err == nil {
		t.Fatalf("Expected error - but got none\n")
	}
	// Check error is 403
	if !strings.HasPrefix(err.Error(), "403 ") {
		t.Fatalf("Expected error 403 - but got %s\n", err.Error())
	}

}
