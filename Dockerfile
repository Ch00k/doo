FROM alpine:3.13

COPY ./doo /doo

ENTRYPOINT ["/doo"]
