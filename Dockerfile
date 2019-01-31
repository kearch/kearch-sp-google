FROM golang:1.10 AS build
WORKDIR /go/src
COPY go ./go
COPY main.go .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o kearch-sp-google .

FROM scratch AS runtime
COPY --from=build /go/src/kearch-sp-google ./
COPY en_default_dict.txt ./
EXPOSE 32500/tcp
ENTRYPOINT ["./kearch-sp-google"]
