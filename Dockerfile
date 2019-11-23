# build stage
FROM golang:alpine AS build
RUN apk --no-cache add build-base git bzr mercurial gcc openrc
ADD . /src/github.com/PumpkinSeed/invoker
RUN cd /src/github.com/PumpkinSeed/invoker/server/cmd && go build -o invoker-http

# final stage
FROM alpine
COPY --from=build /src/github.com/PumpkinSeed/invoker/server/cmd/invoker-http /usr/local/bin
COPY invoker.init /etc/conf.d/
RUN chmod +x /usr/local/bin/invoker-http
RUN rc-service invoker start
EXPOSE 3000

ENTRYPOINT invoker-http