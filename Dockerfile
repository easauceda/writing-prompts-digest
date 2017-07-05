FROM golang as build
WORKDIR /build
ADD main.go template.html ./
RUN go get -d ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app  .

FROM busybox
COPY --from=build /build/app /
ADD template.html .
ADD ca-certificates.crt /etc/ssl/certs/
CMD ["/app"]
