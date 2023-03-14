FROM alpine:latest

RUN mkdir /app

COPY discord /app

CMD [ "/app/discord" ]