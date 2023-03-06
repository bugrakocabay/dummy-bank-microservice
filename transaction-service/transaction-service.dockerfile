FROM alpine:latest
RUN mkdir /app
COPY transactionServiceApp /app
CMD [ "/app/transactionServiceApp" ]