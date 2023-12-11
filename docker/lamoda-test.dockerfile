FROM golang:latest

WORKDIR /lamoda-test

COPY . .

COPY .env .
RUN make build

RUN chmod ugo+x .bin/lamoda-test

EXPOSE 4000

CMD .bin/lamoda-test
