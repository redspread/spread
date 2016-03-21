#!/bin/bash

# setup kubectl config
echo "setting up kubectl config"
mkdir -p ~/.kube
>~/.kube/config cat <<EOF
kind: Config
apiVersion: v1
current-context: gitlab
contexts:
- name: gitlab
  context:
    cluster: default
    user: default
clusters:
- name: default
  cluster:
    server: $KUBECFG_SERVER
    api-version: $KUBECFG_API_VERSION
    insecure-skip-tls-verify: $KUBECFG_INSECURE_SKIP_TLS_VERIFY
    certificate-authority: $KUBECFG_CERTIFICATE_AUTHORITY
users:
- name: default
  user:
    client-certificate: $KUBECFG_CLIENT_CERTIFICATE
    client-certificate-data: $KUBECFG_CLIENT_CERTIFICATE_DATA
    client-key: $KUBECFG_CLIENT_KEY
    token: $KUBECFG_TOKEN
    username: $KUBECFG_USERNAME
    password: $KUBECFG_PASSWORD
EOF

spread deploy $DEPLOY_DIR

