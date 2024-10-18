FROM golang:1.23 AS build
WORKDIR /go/src

# Copy the Go module files and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY go ./go
COPY main.go .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -a -installsuffix cgo -o swagger .

FROM scratch AS runtime
COPY --from=build /go/src/swagger ./
EXPOSE 8080/tcp
ENTRYPOINT ["./swagger"]
