.DEFAULT_GOAL := everything

dependencies:
	@echo Downloading Dependencies
	@go get ./...

build: dependencies
	@echo Compiling Apps
	@go install github.com/riomhaire/lightauth
	@go install github.com/riomhaire/lightauth/lightauthsession
	@go install github.com/riomhaire/lightauth/lightauthuser


test:
	@echo Running Unit Tests
	@go test github.com/riomhaire/lightauth/services

profile:
	@echo Profiling Code
	@go test -coverprofile coverage.out github.com/riomhaire/lightauth/services
	@go tool cover -html=coverage.out -o coverage.html
	@rm coverage.out

clean:
	@echo Cleaning
	@go clean

everything: clean build profile  
	@echo Done
