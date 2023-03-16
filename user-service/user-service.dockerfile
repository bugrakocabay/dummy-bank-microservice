FROM alpine:latest
RUN mkdir /app
COPY userServiceApp /app
COPY config.json /app
CMD [ "/app/userServiceApp" ]