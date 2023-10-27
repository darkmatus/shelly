FROM golang:1.21 AS shelly-dev

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
RUN go install github.com/boumenot/gocover-cobertura@v1.2.0

WORKDIR /opt/app/api
CMD ["air"]