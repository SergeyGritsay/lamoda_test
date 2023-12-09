FROM golang:latest

WORKDIR /lamoda-test

COPY . .

RUN make build

RUN chmod ugo+x .bin/lamoda-test

CMD .bin/lamoda-test
