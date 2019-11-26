FROM golang:latest

LABEL maintainer="Olakunle Odegbaro <oodegbaro@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main server/server.go

EXPOSE 9000

CMD ["./main"]