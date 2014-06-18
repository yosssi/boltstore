protoc:
	protoc --gogo_out=. shared/protobuf/*.proto
bench:
	CREATE_BENCHDATA=true go test -bench . -benchmem -cover ./...
