FROM alpine:latest
RUN mkdir /app
COPY reportServiceApp /app
COPY config.json /app
CMD [ "/app/reportServiceApp" ]