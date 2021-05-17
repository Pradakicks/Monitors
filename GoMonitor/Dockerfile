FROM golang:1.16

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .
COPY . .
EXPOSE 7243
CMD ["/app/main"]