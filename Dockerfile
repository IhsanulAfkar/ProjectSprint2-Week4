FROM golang:1.22


WORKDIR /Week4

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /belimang
EXPOSE 8080
CMD ["/belimang"]