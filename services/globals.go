package services

import (
	"github.com/thoas/stats"
)

// Version
var Version string = "0.7"

// Application Parameters
var SessionSecret *string
var SessionPeriod *int
var UserFile *string // Known users

// Contains the global variables used
var StatsMiddleware = stats.New()

// User stuff
var users = make(map[string]User)

// Session stuff
var ApplicationSessionService = new(SessionService)
