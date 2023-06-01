FROM traefik:3.0

COPY traefik.yaml /etc/traefik/traefik.yaml

COPY conf.yaml /etc/traefik/conf.yaml

RUN apk --no-cache add curl

CMD ["traefik"]
