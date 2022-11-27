# Собираем в гошке
FROM golang:1.19 as build
LABEL ORGANIZATION="PalPalych"
LABEL SERVICE="banner_rotation_int_tests"
LABEL MAINTAINERS="PalPalych@bk.ru"

ENV CODE_DIR /go/src/
WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
COPY ./integration_tests/. ./integration_tests/
CMD go test -tags integration ./integration_tests/...
