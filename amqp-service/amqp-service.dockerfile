FROM alpine:latest
RUN mkdir /app
COPY amqpServiceApp /app
CMD [ "/app/amqpServiceApp" ]