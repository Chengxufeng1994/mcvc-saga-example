.PHONY: proto

proto:
	@echo "Generating protocol buffers"
	rm -rf pb/
	mkdir -p pb/
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto