FROM golang
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY src/ /app/src
COPY config/ /app/config
WORKDIR /app/src