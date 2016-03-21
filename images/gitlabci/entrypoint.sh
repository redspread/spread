#!/bin/sh

# if config file does not exist produce one with env vars
if [ ! -f ~/.kube/config ]; then
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
    certificate-authority-data: $KUBECFG_CERTIFICATE_AUTHORITY_DATA

users:
- name: default
  user:
    client-certificate: $KUBECFG_CLIENT_CERTIFICATE
    client-certificate-data: $KUBECFG_CLIENT_CERTIFICATE_DATA
    client-key: $KUBECFG_CLIENT_KEY
    client-key-data: $KUBECFG_CLIENT_KEY_DATA
    token: $KUBECFG_TOKEN
    username: $KUBECFG_USERNAME
    password: $KUBECFG_PASSWORD
EOF
fi

if [ ! -z "$DEPLOY_DIR" ]; then
    echo "using DEPLOY_DIR..."
    spread deploy $DEPLOY_DIR
    exit
fi

echo "Using project dir..."
spread deploy $CI_PROJECT_DIR