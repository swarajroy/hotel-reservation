build:
	@go build -o bin/api

run: build
	@./bin/api
docker:
	echo "building docker file"
	@docker build -t hotel_reservation_api:latest .


seed:
	@go run scripts/seed.go

test:
	@go test -v -count=1 ./... 