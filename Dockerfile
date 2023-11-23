FROM alpine:3.18

MAINTAINER Dmitry Mozzherin

ENV LAST_FULL_REBUILD 2023-11-22

WORKDIR /bin

COPY ./bhlquest /bin

ENTRYPOINT [ "bhlquest", "serve" ]

CMD ["-p", "8555"]
