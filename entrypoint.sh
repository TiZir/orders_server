wait-for "${PG_HOST}:${PG_PORT}" -- "$@"
go build -o pub ./cmd/main.go
chmod 777 ./pub
go build -o main main.go
chmod 777 ./main
./main &
./pub &