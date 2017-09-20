package services

import (
	"testing"
)

func TestStatisticsService(t *testing.T) {
	service := new(StatsService)
	args := StatsArgs{}
	result := StatsReply{}
	result.Statistics = StatsMiddleware.Data()

	err := service.Status(nil, &args, &result)
	if err != nil {
		t.Fatalf("Expected a result but got %s\n", err.Error())
	}
}
