ARG APP_NAME=htmx-dashboard

FROM golang:1.21-alpine as build

ARG APP_NAME
WORKDIR /app

RUN apk --no-cache add git alpine-sdk
COPY go.* ./
RUN go mod download
COPY *.go ./
RUN GO111MODULE=on go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./$APP_NAME

FROM scratch

ARG APP_NAME

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/$APP_NAME /$APP_NAME
COPY www/html/* /www/html/
COPY www/css/* /www/css/
COPY Favorites.json /Favorites.json

EXPOSE 8182
ENTRYPOINT ["/htmx-dashboard"]
