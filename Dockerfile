# build stage
FROM golang:alpine AS build
RUN apk --no-cache add build-base git bzr mercurial gcc
RUN useradd
ADD . /src/github.com/PumpkinSeed/invoker
RUN cd /src/github.com/PumpkinSeed/invoker/server/cmd && go build -o invoker-http

# final stage
FROM alpine
COPY --from=build /src/github.com/PumpkinSeed/invoker/server/cmd/invoker-http /usr/local/bin
RUN chmod +x /usr/local/bin/invoker-http
EXPOSE 3000

ENTRYPOINT invoker-http