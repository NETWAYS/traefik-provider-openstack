FROM golang:alpine as build

ENV CGO_ENABLED=0
COPY ./ /build
WORKDIR /build
RUN go get -v ./...
RUN go build -o traefik-provider-openstack .

FROM alpine
COPY --from=build /build/traefik-provider-openstack /usr/local/bin/
ENTRYPOINT ["traefik-provider-openstack"]
USER nobody
