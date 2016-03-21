FROM gliderlabs/alpine:3.3

RUN apk-install bash

WORKDIR /opt/spread-gitlab

ADD spread-linux-static /usr/local/bin/spread
ADD entrypoint.sh ./entrypoint.sh

ENV KUBECFG_INSECURE_SKIP_TLS_VERIFY="false"
ENV DEPLOY_DIR="."

ENTRYPOINT ["./entrypoint.sh"]