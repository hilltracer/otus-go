//go:generate protoc -I=../../api --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ../../api/EventService.proto
package pb
