FROM golang:alpine

ARG app_name=app
RUN mkdir -p /opt/${app_name}

WORKDIR /opt/${app_name}

COPY . .

RUN go mod download && \
unset http_proxy && \
unset https_proxy && \
CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o go_app cmd/ws/main.go

FROM alpine

ARG app_name=app
ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

COPY --from=0 /opt/${app_name}/go_app /app/go_app
COPY --from=0 /opt/${app_name}/config.yaml /app/config.yaml

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

CMD ["/app/go_app"]