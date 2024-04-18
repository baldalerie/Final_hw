run:
	docker-compose -f docker-compose.yml up -d
	go run .

stop:
	docker-compose -f docker-compose.yml down


test:
	docker-compose -f docker-compose.yml up -d
	go test -v ./...
	docker-compose -f docker-compose.yml down


test-db-up:
	docker-compose -f docker-compose.yml up -d

test-db-down:
	docker-compose -f docker-compose.yml down

clean:
	go clean