FROM alpine

WORKDIR /workspace/ticket

COPY ticket .

ADD config.yml .
ADD log.json .
ADD logs .

EXPOSE 8041

CMD ["./ticket"]