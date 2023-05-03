FROM golang:1.20.3-bullseye as builder

ENV WORKDIR /home/appuser/app
ENV BUILDPATH $WORKDIR/customer-manager

RUN useradd -s /bin/bash -U -u 10001 appuser -m
USER appuser
WORKDIR $WORKDIR
COPY --chown=appuser:appuser . $WORKDIR

RUN go build -o $BUILDPATH $WORKDIR/cmd/server/server.go
RUN go build -o ${BUILDPATH}-migartions $WORKDIR/cmd/migrate/migrate.go

FROM golang:1.20.3-bullseye as customer-manager

ENV WORKDIR /home/appuser/app
ENV BUILDPATH $WORKDIR/customer-manager

WORKDIR ${WORKDIR} 

RUN useradd -s /bin/bash -U -u 10001 appuser -m
USER appuser
COPY --from=builder $BUILDPATH /usr/local/bin/customer-manager
COPY --from=builder ${BUILDPATH}-migartions /usr/local/bin/customer-manager-migrations

CMD customer-manager
