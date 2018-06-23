#Build stage
FROM golang:1.10-alpine3.7 AS build-env

#Build deps
RUN apk --no-cache add build-base git

#Setup repo
COPY . ${GOPATH}/src/git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator
WORKDIR ${GOPATH}/src/git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator

RUN make docker-binary

FROM alpine:3.7
LABEL maintainer="info@jonasfranz.software"

COPY --from=build-env /go/src/git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/gitea-github-migrator /usr/local/bin/gitea-github-migrator

ENTRYPOINT ["/usr/local/bin/gitea-github-migrator"]
