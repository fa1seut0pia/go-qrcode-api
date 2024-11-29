FROM golang:alpine AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
#RUN export GOPROXY=https://goproxy.io,direct && go mod download
RUN go mod download
COPY main.go .
RUN go build -o go-qrcode-api

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/go-qrcode-api .
COPY SarasaFixedSC-Regular.ttf .
EXPOSE 7688
CMD ["/app/go-qrcode-api"]