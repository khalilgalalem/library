FROM golang:1.10 AS build
WORKDIR /go/src

# Download the specific compatible version of Gorilla Mux
RUN go get -d -v github.com/gorilla/mux && \
    cd /go/src/github.com/gorilla/mux && \
    git checkout v1.7.4

COPY go ./go
COPY main.go .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o swagger .

FROM scratch AS runtime
COPY --from=build /go/src/swagger ./
EXPOSE 8080/tcp
ENTRYPOINT ["./swagger"]
