FROM golang as build
WORKDIR /build
ADD main.go template.html ./
RUN go get -d ./...
RUN curl -o ca-certificates.crt https://raw.githubusercontent.com/bagder/ca-bundle/master/ca-bundle.crt
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app  .

FROM scratch
COPY --from=build /build/app /
COPY --from=build /build/ca-certificates.crt /etc/ssl/certs/
ADD template.html .
CMD ["/app"]
