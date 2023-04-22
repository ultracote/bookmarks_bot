FROM golang:1.20-alpine AS build

WORKDIR /app

RUN apk --no-cache add ca-certificates

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app cmd/main.go


FROM scratch

WORKDIR /app

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/app .

CMD ["./app"]