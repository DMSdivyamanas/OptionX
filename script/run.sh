# Download the project dependencies
go mod download

# Build the project
go build -o server ./cmd/server/main.go

# Run the server
./server