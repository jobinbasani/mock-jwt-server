FROM public.ecr.aws/bitnami/golang:1.15.5 as builder

ENV GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . /build/src

RUN cd src && \
    go build -o ../out/mock-jwt-server

FROM alpine
COPY --from=builder /build/out /jwt

ENTRYPOINT ["/jwt/mock-jwt-server"]

EXPOSE 8989