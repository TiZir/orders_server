wait-for "${PG_HOST}:${PG_PORT}" -- "$@"
go build -o main main.go
go build -o pub ./cmd/main.go
./pub &
./main 