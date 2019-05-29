FROM golang:1.12 AS build
MAINTAINER Micah Huber <micah@reactiveops.com>
WORKDIR /go/src/github.com/reactiveops/dd-manager
ADD . /go/src/github.com/reactiveops/dd-manager

RUN GO111MODULE=on GOOS=linux GOARCH=amd64 go build


FROM gcr.io/distroless/base
COPY --from=build /go/src/github.com/reactiveops/dd-manager/dd-manager /dd-manager
COPY --from=build /go/src/github.com/reactiveops/dd-manager/conf.yml /conf.yml
ENTRYPOINT ["/dd-manager"]