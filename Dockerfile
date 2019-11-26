# Наследуемся от alpine3.7
FROM golang:1.10.1-alpine3.7 AS builder

# ENV
ENV APPDIR $GOPATH/src/github.com/eekrupin/offersStore

# FS
RUN mkdir -p ${APPDIR}
WORKDIR ${APPDIR}

COPY api api
COPY config config
COPY db db
COPY models models
COPY modules modules
COPY queries queries
COPY services services
COPY vendor vendor
COPY main.go .

RUN ls -lah

RUN go build -ldflags "-s -w" -o /build/app

RUN ls -lah /build

FROM alpine:3.7

COPY --from=builder /build/app /app
#COPY queries queries

CMD ["/app"]

EXPOSE 80
