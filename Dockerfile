FROM alpine:3.13

COPY pluto-vs-relay.linux64 /usr/local/bin/pluto-vs-relay
USER nobody
CMD /usr/local/bin/pluto-vs-relay