BUILD_DIR=$(PWD)/bin
profile:
	echo "Building profile service"
	cd src/golang/profile-service && GOOS=linux GOARCH=386 go build -o ${BUILD_DIR}/profile-service main.go

profile-test:
	echo "Running unit tests for profile service"
	cd src/golang/profile-service && go test -short ./...

profile-test-integration:
	echo "Running integration tests for profile service"
	cd src/golang/profile-service \
		&& docker build . -t profile-service \
		&& mkdir -p test_data \
		&& cp ../../../schema/players.sql test_data/schema.sql \
		&& go test --tags=integration ./...

matchmaking:
	echo "Building matchmaking service"
	cd src/golang/matchmaking-service && GOOS=linux GOARCH=386 go build -o ${BUILD_DIR}/matchmaking-service main.go

matchmaking-test:
	echo "Testing unit tests for matchmaking service"
	cd src/golang/matchmaking-service && go test -short ./...

matchmaking-test-integration:
	echo "Running integration tests for matchmaking service"
	cd src/golang/matchmaking-service \
		&& docker build . -t matchmaking-service \
		&& mkdir -p test_data \
		&& cp ../../../schema/players.sql test_data/schema.sql \
		&& go test --tags=integration ./...

item:
	echo "Building item service"
	cd src/golang/item-service && GOOS=linux GOARCH=386 go build -o ${BUILD_DIR}/item-service main.go

item-test:
	echo "Testing item service"
	cd src/golang/item-service && go test ./...

tradepost:
	echo "Building tradepost service"
	cd src/golang/tradepost-service && GOOS=linux GOARCH=386 go build -o ${BUILD_DIR}/tradepost-service main.go

tradepost-test:
	echo "Testing tradepost service"
	cd src/golang/tradepost-service && go test ./...

test-all-unit: profile-test matchmaking-test item-test tradepost-test

test-all-integration: profile-test-integration matchmaking-test-integration

test-all: test-all-unit test-all-integration

build-all: profile matchmaking item tradepost

clean:
	echo "Running cleanup"
	rm bin/*
	docker rmi -f profile-service
