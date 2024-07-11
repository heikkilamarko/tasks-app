FROM golang:1.23-rc AS build
ARG GOARCH=amd64
COPY . .
ENV GOPATH=""
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=$GOARCH
RUN go build -trimpath -a -ldflags="-w -s" ./cmd/tasks-app

FROM gcr.io/distroless/static
COPY --from=build /go/tasks-app /tasks-app
ENTRYPOINT ["/tasks-app"]
