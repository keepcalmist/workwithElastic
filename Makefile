biuld:
rm -rf build && mkdir build && build -o build/main -v ./cmd/api/
run:
go run cmd/api/main.go