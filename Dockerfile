FROM golang:alpine

COPY . ./build
WORKDIR /build
RUN GOBIN=`pwd` CGO_ENABLED=0 go install go.uber.org/sally@latest

FROM scratch
COPY --from=0 /build/sally sally
COPY sally.yaml /
EXPOSE 8080
WORKDIR /
CMD ["/sally"]
