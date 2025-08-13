FROM node:lts AS build-ui
WORKDIR /app
COPY internal/modules/ui/web/package*.json .
RUN npm ci
COPY internal/modules/ui/web/ .
RUN npm run build

FROM golang:1.25 AS build
COPY . .
COPY --from=build-ui /app/dist/ ./internal/modules/ui/
ENV GOPATH=""
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -trimpath -a -ldflags="-w -s" ./cmd/tasks-app

FROM gcr.io/distroless/static
COPY --from=build /go/tasks-app /tasks-app
ENTRYPOINT ["/tasks-app"]
