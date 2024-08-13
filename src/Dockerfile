FROM node:20-alpine AS build-ui
WORKDIR /app
COPY internal/modules/ui/web/package*.json .
RUN npm ci
COPY internal/modules/ui/web/ .
RUN npm run build

FROM golang:1.23-rc AS build
ARG GOARCH=amd64
COPY . .
COPY --from=build-ui /app/dist/ ./internal/modules/ui/static/
COPY --from=build-ui /app/dist/index.html ./internal/modules/ui/templates/index.html
ENV GOPATH=""
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=$GOARCH
RUN go build -trimpath -a -ldflags="-w -s" ./cmd/tasks-app

FROM gcr.io/distroless/static
COPY --from=build /go/tasks-app /tasks-app
ENTRYPOINT ["/tasks-app"]
