FROM scratch

ENV AUSTIN_SERVER_HOST 127.0.0.1
ENV AUSTIN_SERVER_PORT 8084

EXPOSE $AUSTIN_SERVER_PORT

COPY bin/linux-amd64/austin /
COPY config/config.toml /config/

CMD ["/austin"]
