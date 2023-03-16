FROM alpine:latest
RUN mkdir /app
COPY loggerServiceApp /app
COPY config.json /app
CMD [ "/app/loggerServiceApp" ]