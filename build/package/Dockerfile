FROM golang:alpine AS build
WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg

RUN CGO_ENABLED=0 go build -v -o /usr/src/app/app -ldflags '-extldflags "-static"' cmd/analyse/main.go

FROM scratch
WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/src/app/app ./

CMD ["/app"]