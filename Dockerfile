# Minimize container size using multi stage build
FROM golang:1.21-alpine as builder

WORKDIR /app

ENV CGO_ENABLED=false\
    GOOS=linux\
    GOARCH=amd64

# Caching
COPY go.mod go.sum* ./

RUN go mod download

COPY . .

RUN go build -o dexus cmd/main.go


FROM scratch as final

COPY --from=builder /app/dexus .

ENTRYPOINT [ "./dexus" ]
