FROM golang:1.21-bullseye as builder

WORKDIR /app
EXPOSE 8080

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 go build -o bin/

FROM gcr.io/distroless/static-debian11

COPY --from=builder /app/bin/metrix .

CMD ["./metrix", "serve"]
