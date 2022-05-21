FROM golang:1.18

WORKDIR /usr/local/mafia_gql_server
COPY db db
COPY graph graph
COPY utils utils
COPY *.go ./
COPY go.mod .
COPY go.sum .

RUN go build .

CMD ["./MafiaGQL_server"]
