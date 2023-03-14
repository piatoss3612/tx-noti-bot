FROM alpine:latest

RUN mkdir /app

COPY auth /app

CMD [ "/app/auth" ]