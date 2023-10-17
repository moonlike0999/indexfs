FROM caddy:2.7.5-builder AS builder

RUN xcaddy build \
    --with github.com/moonlike0999/indexfs/caddy \
    --with github.com/caddy-dns/digitalocean \
    --with github.com/caddyserver/json5-adapter

FROM caddy:2.7.5

COPY --from=builder /usr/bin/caddy /usr/bin/caddy

RUN mkdir /cfg /fs
RUN echo "{}" >> /cfg/config.json5

CMD ["caddy", "run", "--config", "/cfg/config.json5", "--adapter", "json5"]