FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /backend-2 .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /backend-2 /backend-2
EXPOSE 8080
ENTRYPOINT ["/backend-2"]
