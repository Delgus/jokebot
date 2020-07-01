FROM alpine:3.11
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY bot bot
EXPOSE 80
CMD ["./bot"]