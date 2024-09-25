FROM golang:1.23-bullseye as build
LABEL stage="builder"

WORKDIR /go/src/app
ADD . /go/src/app

RUN go mod download

RUN go build -o /go/bin/autovoter /go/src/app


FROM gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=build /go/bin/autovoter ./app

CMD ["./app"]
