FROM golang
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go get -u github.com/swaggo/swag/cmd/swag
COPY src/ /app/src
COPY config/ /app/config
WORKDIR /app/src
RUN swag init