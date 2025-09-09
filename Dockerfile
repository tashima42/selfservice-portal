FROM golang:1.25 AS builder

ARG VERSION=dev

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /selfserviceportal -ldflags "-w -X main.Version=${VERSION}" .

FROM alpine

WORKDIR /

COPY --from=builder /selfserviceportal /selfserviceportal

USER 1001:1001

ENTRYPOINT ["/selfserviceportal"]
