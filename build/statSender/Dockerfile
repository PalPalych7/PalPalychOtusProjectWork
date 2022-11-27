# Собираем в гошке
FROM golang:1.19 as build

ENV BIN_FILE /opt/banner_stat_sender/banner_stat_sender-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/statSender/*

FROM alpine:3.9

LABEL ORGANIZATION="PalPalych"
LABEL SERVICE="banner_stat_sender"
LABEL MAINTAINERS="PalPalych@bk.ru"

ENV BIN_FILE "/opt/banner_stat_sender/banner_stat_sender-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/banner_stat_sender/banner_stat_sender-app/config.toml
COPY ./configs/statSenderConfig.toml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}


