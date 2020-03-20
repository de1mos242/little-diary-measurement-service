FROM golang as build-env
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY src/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /app/measurement-service

FROM scratch
COPY --from=build-env /app/measurement-service /app/measurement-service
EXPOSE 8080
ENTRYPOINT ["/app/measurement-service"]