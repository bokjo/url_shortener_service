FROM golang:latest as build

RUN mkdir /shortener
WORKDIR /shortener
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . . 

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/shortenerapi

FROM scratch
COPY --from=build /go/bin/shortenerapi /go/bin/shortenerapi
ENTRYPOINT ["/go/bin/shortenerapi"]