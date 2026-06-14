FROM golang:1.26.3 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/down-detector ./cmd

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=build /out/down-detector /app/down-detector
COPY web /app/web

EXPOSE 5436

CMD ["/app/down-detector"]
