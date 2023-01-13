up:
	docker-compose up -d
down:
	docker-compose down
run:
	go run main.go
runclient:
	go run client/main.go
proto:
	mkdir pb
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		proto/*.proto
callrpc:
	evans --host localhost --port 8080 -r repl
