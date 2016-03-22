FROM docker:git

ADD spread-linux-static /usr/local/bin/spread
ADD entrypoint.sh /opt/spread-gitlab/entrypoint.sh

ENV KUBECFG_INSECURE_SKIP_TLS_VERIFY="false"

ENTRYPOINT ["/opt/spread-gitlab/entrypoint.sh"]