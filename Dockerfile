FROM golang as build-env
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go get -u github.com/swaggo/swag/cmd/swag
COPY src/ /app/src
WORKDIR /app/src
RUN swag init
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /app/measurement-service

FROM scratch
COPY --from=build-env /app/measurement-service /app/measurement-service
COPY config/ /app/config
WORKDIR /app
ENV GIN_MODE release
EXPOSE 8080
ENTRYPOINT ["/app/measurement-service"]