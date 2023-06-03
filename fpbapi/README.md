### Update proto
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative feed-service.proto
```

### Update mocks
```
mockery --name FeedServiceClient
```
