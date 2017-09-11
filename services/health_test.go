package services

import (
	"testing"
)

func TestHealthService(t *testing.T) {
	service := new(HealthService)
	args := HealthArgs{}
	result := HealthReply{}

	err := service.Status(nil, &args, &result)
	if err != nil {
		t.Fatalf("Expected a result but got %s\n", err.Error())
	}
}
