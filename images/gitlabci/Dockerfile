FROM gliderlabs/alpine:3.3

RUN apk-install bash

WORKDIR /opt/spread-gitlab

ADD spread-linux /usr/local/bin/spread
ADD entrypoint.sh ./entrypoint.sh

ENV KUBECFG_INSECURE_SKIP_TLS_VERIFY="false"

ENTRYPOINT ["entrypoint.sh"]