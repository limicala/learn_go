docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.30.0-alpine \
sh -c "golangci-lint run" -v
