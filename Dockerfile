# pakai image golang
FROM golang:1.23-alpine AS builder

# install tzdata di build stage
RUN apk add --no-cache tzdata

WORKDIR /app

# copy dependency
COPY go.mod go.sum ./
RUN go mod download


# copy source
COPY . .

# build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/web

# stage runtime (lebih kecil)
FROM alpine:latest  

# install tzdata juga di runtime stage
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/server .

# expose port
EXPOSE 3000

CMD ["./server"]
