FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /agent ./cmd/agent/main.go

FROM scratch
WORKDIR /
COPY --from=builder /api /api
# In a real environment, the agent would be deployed separately
EXPOSE 8080
ENTRYPOINT ["/api"]
