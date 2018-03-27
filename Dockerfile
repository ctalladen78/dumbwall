FROM golang:1.10.0-stretch

WORKDIR /go/src/github.com/maksadbek/dumbwall
COPY . /go/src/github.com/maksadbek/dumbwall
RUN make httpd

ENV DSN postgres://postgres:mysecretpassword@localhost:5432/dumbwall?sslmode=disable

VOLUME /etc/dumbwall

EXPOSE 80

CMD ["httpd"]
