go run main.go --restport 10000 --comuport 20000 --addr localhost:20001 --addr localhost:20002
go run main.go --restport 10001 --comuport 20001 --addr localhost:20000 --addr localhost:20002
go run main.go --restport 10002 --comuport 20002 --addr localhost:20001 --addr localhost:20000
