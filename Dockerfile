FROM golang:1.20.1

WORKDIR /harmony

COPY cmd cmd
COPY internal internal
COPY go.mod go.mod
COPY go.sum go.sum

RUN go build -o /bin/harmony ./cmd/harmony

ENTRYPOINT ["harmony"]