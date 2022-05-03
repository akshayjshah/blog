# syntax=docker/dockerfile:1
FROM amd64/golang:alpine AS build
RUN apk update
RUN apk upgrade
RUN apk add git make bash
ADD . /src
WORKDIR /src
RUN rm -f tmp/serve && make tmp/serve

FROM amd64/alpine
WORKDIR /app
COPY --from=build /src/tmp/serve /app/serve
RUN apk update
RUN apk upgrade
RUN apk add ca-certificates tzdata
CMD [ "./serve" ]
