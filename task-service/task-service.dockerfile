FROM alpine:latest

RUN mkdri /app

COPY taskApp /app

CMD ["/app/taskApp"]