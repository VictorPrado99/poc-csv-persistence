FROM golang:1.18

ENV PERSISTENCE_HOME="/.csv_persistence"

COPY ./config.json /.csv_persistence/config.json

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . /usr/src/app
RUN go build -v -o /usr/local/bin/app

EXPOSE 9001

CMD ["app"]
