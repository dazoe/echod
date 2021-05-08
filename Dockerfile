FROM golang:1.16-alpine as builder

RUN apk -u add git;

WORKDIR /src
COPY . .

RUN go build -o /echod .

FROM alpine:3.13

COPY --from=builder /echod /

ENTRYPOINT [ "/echod" ]
CMD ["-l", ":7"]