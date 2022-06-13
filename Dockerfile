ARG MIGRATE_VERSION=v4.13.0
ARG CREATED
ARG REVISION
ARG TITLE
ARG VERSION

FROM migrate/migrate:${MIGRATE_VERSION} as migrate_image
FROM quay.io/elasticpath/lang-go:1.16-e67ebccd-35388

ENV GOPATH=/go
WORKDIR /app
COPY . .

RUN apk add --no-cache bash=~5 git

RUN go build -o build/app gitlab.elasticpath.com/commerce-cloud/personal-data.svc/cmd/app


# MetaData
LABEL org.opencontainers.image.vendor="Elastic Path Software Inc."
LABEL org.opencontainers.image.source="https://gitlab.elasticpath.com/commerce-cloud/personal-data.svc"
LABEL org.opencontainers.image.title=$TITLE
LABEL org.opencontainers.image.description="This Docker image contains the Personal Data service of Elastic Path Commerce Cloud application."
LABEL org.opencontainers.image.created=$CREATED
LABEL org.opencontainers.image.version=$VERSION
LABEL org.opencontainers.image.revision=$REVISION

RUN mv /app/build/app /personal-data

# migrations
RUN mkdir -p ./db
COPY ./db /db
COPY ./entrypoint.sh /
COPY --from=migrate_image /usr/local/bin/migrate /migrate
RUN chmod +x /migrate
RUN chmod +x /entrypoint.sh


# This seems like something we should set
# https://github.com/gin-gonic/gin/blob/master/mode.go#L20
ENV GIN_MODE=release

ENTRYPOINT ["/entrypoint.sh"]
