FROM golang:alpine AS builder

RUN apk add --no-cache --update git build-base
ENV CGO_ENABLED=1
ARG VERSION=0.0.0-dev
ARG COMMIT=none
ARG DATE=unknown
WORKDIR /app
COPY ./src/go.mod ./src/go.sum ./
RUN go mod download

COPY ./src .
RUN go build -ldflags "-s -w -X github.com/assimon/luuu/config.BuildVersion=${VERSION} -X github.com/assimon/luuu/config.BuildCommit=${COMMIT} -X github.com/assimon/luuu/config.BuildDate=${DATE}"  \
	-o epusdt .

FROM alpine:latest AS runner
ENV TZ=Asia/Shanghai
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/static /app/static
COPY --from=builder /app/static /static
COPY --from=builder /app/epusdt .
VOLUME /app/conf

ENTRYPOINT ["./epusdt" ,"http","start"]
