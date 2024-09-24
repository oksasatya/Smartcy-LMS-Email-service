PROTOC_GEN_GO := protoc --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative

PROTOS := proto/email.proto


all: generate

# Generate Go files from proto
generate:
	$(PROTOC_GEN_GO) $(PROTOS)

# Clean generated files
clean:
	rm -rf pb/*.pb.go

.PHONY: all generate clean