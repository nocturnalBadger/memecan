FROM golang:alpine

WORKDIR /memecan
COPY . /memecan

RUN go build -o memecan .

CMD ./memecan
