FROM scratch

ENV AUSTIN_SERVER_HOST 127.0.0.1
ENV AUSTIN_SERVER_PORT 8084

ENV AUSTIN_STORAGE_HOST cassandra.k8s.gromnsk.ru
ENV AUSTIN_CONSUL_HOSTPORT consul.k8s.gromnsk.ru:8500

EXPOSE $AUSTIN_SERVER_PORT

COPY bin/linux-amd64/austin /
COPY config/config.toml /config/

CMD ["/austin"]
