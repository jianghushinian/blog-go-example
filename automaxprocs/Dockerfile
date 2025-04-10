FROM alpine:latest
LABEL authors="jianghushinian"

WORKDIR /app

COPY . .

# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gomaxprocs main.go
# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o automaxprocs/automaxprocs automaxprocs/main.go

# docker build -t automaxprocs:latest .

# docker run --cpus=2 -it \
#   -v $(pwd):/app -w /app alpine \
#   ./gomaxprocs

# docker run --cpus=2 -it \
#   -v $(pwd):/app -w /app alpine \
#   ./automaxprocs/automaxprocs
