# build
FROM golang:1.10.3-alpine3.7 AS build-env
ENV REPO_URI=github.com/rayou/theperfectgiftcard-api
ADD . $GOPATH/src/$REPO_URI/
RUN cd $GOPATH/src/$REPO_URI/ && go build -o app && mv ./app /usr/bin/app

# final
FROM alpine
MAINTAINER Ray Ou <yuhung.ou@live.com>
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=build-env /usr/bin/app /usr/local/bin/app
ENV PORT=8000
EXPOSE $PORT
CMD app
