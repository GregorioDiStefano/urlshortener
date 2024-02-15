FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.sh ./
COPY *.go ./


RUN CGO_ENABLED=0 GOOS=linux go build -o /urlshort

EXPOSE 8888

RUN chmod +x ./docker-env-set.sh

ENTRYPOINT ["bash", "-c", "source /app/docker-env-set.sh && /urlshort"]