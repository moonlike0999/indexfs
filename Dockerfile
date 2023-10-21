FROM caddy:2.7.5-builder AS builder

RUN xcaddy build \
    --with github.com/moonlike0999/indexfs/caddyfs \
    --with github.com/caddy-dns/digitalocean \
    --with github.com/caddyserver/json5-adapter

FROM caddy:2.7.5

COPY --from=builder /usr/bin/caddy /usr/bin/caddy

RUN mkdir /fs
RUN echo "{}" >> /config/config.json5
RUN rm /config/Caddyfile

CMD ["caddy", "run", "--config", "/config/config.json5", "--adapter", "json5"]