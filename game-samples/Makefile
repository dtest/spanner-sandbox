BUILD_DIR=$(PWD)/bin
profile:
	echo "Building profile service"
	cd src/golang/profile-service && GOOS=linux GOARCH=386 go build -o ${BUILD_DIR}/profile-service main.go

profile-test:
	echo "Testing profile service"
	cd src/golang/profile-service && go test ./...

matchmaking:
	echo "Building matchmaking service"
	cd src/golang/matchmaking-service && GOOS=linux GOARCH=386 go build -o ${BUILD_DIR}/matchmaking-service main.go

matchmaking-test:
	echo "Testing matchmaking service"
	cd src/golang/matchmaking-service && go test ./...

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

test-all: profile-test matchmaking-test item-test tradepost-test

build-all: profile matchmaking item tradepost

clean:
	echo "Removing binaries"
	rm bin/*
