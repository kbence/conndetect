FROM golang AS build

COPY . /opt/conndetect
WORKDIR /opt/conndetect
RUN make

FROM alpine

RUN apk add --update gcompat && \
    rm -rf /var/cache/apk
COPY --from=build /opt/conndetect/out/conndetect /usr/bin/conndetect

# ENTRYPOINT ["/usr/bin/conndetect"]
# CMD ls /usr/bin
