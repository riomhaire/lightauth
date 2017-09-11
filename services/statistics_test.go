package services

import (
	"testing"

	"github.com/thoas/stats"
)

func TestStatisticsService(t *testing.T) {
	service := new(StatsService)
	args := StatsArgs{}
	result := stats.Data{}

	err := service.Status(nil, &args, &result)
	if err != nil {
		t.Fatalf("Expected a result but got %s\n", err.Error())
	}
}
