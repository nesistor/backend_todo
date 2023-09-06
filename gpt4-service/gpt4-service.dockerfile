FROM alpine:latest

RUN mkdri /app

COPY gpt4App /app

CMD ["/app/gpt4App"]