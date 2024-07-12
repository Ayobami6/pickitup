FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go install github.com/air-verse/air@latest

EXPOSE 2400

CMD [ "air" ]