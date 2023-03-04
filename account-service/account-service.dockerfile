FROM alpine:latest
RUN mkdir /app
COPY accountServiceApp /app
CMD [ "/app/accountServiceApp" ]