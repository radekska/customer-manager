FROM golang:1.20.3-bullseye as builder

ENV WORKDIR /home/appuser/app
ENV BUILDPATH /home/appuser/build
ENV SERVER ${BUILDPATH}/customer-manager
ENV MIGRATIONS ${BUILDPATH}/customer-manager-migrations 

RUN useradd -s /bin/bash -U -u 10001 appuser -m
USER appuser
WORKDIR $WORKDIR
COPY --chown=appuser:appuser . $WORKDIR

RUN go build -o  ${SERVER} ${WORKDIR}/cmd/server/server.go
RUN go build -o  ${MIGRATIONS} ${WORKDIR}/cmd/migrate/migrate.go

FROM golang:1.20.3-bullseye as customer-manager

RUN useradd -s /bin/bash -U -u 10001 appuser
USER appuser
COPY --from=builder ${SERVER} /usr/local/bin/customer-manager
COPY --from=builder ${MIGRATIONS} /usr/local/bin/customer-manager-migrations

CMD customer-manager
