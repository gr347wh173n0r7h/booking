FROM golang:1.15-alpine

RUN apk add git

RUN mkdir /app

RUN cd / && git clone https://github.com/wordnik/swagger-ui.git
ENV SWAGGERDIST /swagger-ui/dist

ADD . /app

WORKDIR /app

RUN go build -o booking cmd/booking/main.go

CMD ["/app/booking"]