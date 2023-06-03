start:
	docker-compose up -d

stop:
	docker-compose down

protogen:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	cd fpbapi && \
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative feed-service.proto

mockgen:
	go install github.com/vektra/mockery/v2@v2.24.0
	go generate ./...

test_with_db:
	docker-compose -f docker-compose.test.yaml up --build --abort-on-container-exit && \
    docker-compose -f docker-compose.test.yaml rm -fsv
