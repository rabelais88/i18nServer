FROM golang:1.16.3-alpine
WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build main.go

# 1. expose is not allowed for heroku app
# 2. $PORT is automatically registered by heroku

CMD i18nServer