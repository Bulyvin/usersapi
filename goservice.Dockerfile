FROM golang:1.13.6 as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build


FROM scratch
COPY --from=builder /app/mod /

ENTRYPOINT ["/mod"]