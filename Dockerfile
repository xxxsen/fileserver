FROM golang:1.22

WORKDIR /build
COPY . ./
RUN CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w' -o file_server ./cmd

FROM alpine:3.12
COPY --from=0 /build/file_server /bin/

ENTRYPOINT [ "/bin/file_server" ]