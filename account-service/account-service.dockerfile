FROM alpine:latest
RUN mkdir /app
COPY accountServiceApp /app
COPY config.json /app
CMD [ "/app/accountServiceApp" ]