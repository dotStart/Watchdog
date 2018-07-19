FROM alpine

RUN apk add --no-cache --update musl-dev && \
    mkdir -p /opt/watchdog/bin && \
    mkdir -p /opt/watchdog/config && \
    addgroup watchdog && \
    adduser -D -g "" -s /bin/false -G watchdog watchdog && \
    chown watchdog:watchdog /opt/watchdog
COPY watchdog /opt/watchdog/bin

VOLUME /opt/watchdog/config
EXPOSE 46624/tcp

USER watchdog
WORKDIR /opt/watchdog
ENTRYPOINT ["/opt/watchdog/bin/watchdog", "-config=/opt/watchdog/config"]
