# syntax=docker/dockerfile:1
ARG GO_VERSION=1.22.5

FROM golang:${GO_VERSION} AS builder

WORKDIR /src

COPY ./go.mod /src/
COPY ./go.sum /src/

RUN go mod download

COPY . /src/

RUN CGO_ENABLED=0 \
    go build -o main

# -------------------------------

FROM scratch

WORKDIR /app

COPY --from=builder /src/main .

EXPOSE 8080

ENTRYPOINT ["./main"]




