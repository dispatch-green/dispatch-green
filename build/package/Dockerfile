FROM golang:1.18 as builder

RUN mkdir /queuebot
ADD . /queuebot
WORKDIR /queuebot
RUN go build -o /queuebot/build/bin/server cmd/queuebot/main.go

# prod container
FROM bullseye

# copy binary from builder container
COPY --from=builder /queuebot/build/bin/server /usr/bin/queuebot

# copy in default config
COPY configs/docker.yml /etc/queuebot/queuebot.yml

RUN mkdir -p /var/www/queuebot

EXPOSE 8080/tcp

ENTRYPOINT [ "/usr/bin/marvin" ] 
