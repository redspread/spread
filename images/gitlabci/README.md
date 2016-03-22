# Gitlab Spread Plugin

This plugin allows you to add `spread` Kubernetes deployment to your Gitlab CI build pipeline.

## Environment Variables

These environment variables can be used to configure the deployment.

### spread

`DEPLOY_DIR`

The directory spread should be deployed from


### Cluster
`KUBECFG_SERVER`

Address of the Kubernetes API Server (https://hostname:port)

`KUBECFG_API_VERSION`

Preferred api version for communicating with the kubernetes cluster (v1, v2, etc)

`KUBECFG_INSECURE_SKIP_TLS_VERIFY`

Disable requirement that connections must pass TLS verification

`KUBECFG_CERTIFICATE_AUTHORITY`

Path to a cert file for the certificate authority.

`KUBECFG_CERTIFICATE_AUTHORITY_DATA`

Certificate data

### User

`KUBECFG_CLIENT_CERTIFICATE`

Path to a client cert file for TLS

`KUBECFG_CLIENT_CERTIFICATE_DATA`

TLS client cert data

`KUBECFG_CLIENT_KEY`

Path to a client key file for TLS.

`KUBECFG_CLIENT_KEY_DATA`

Key data

`KUBECFG_TOKEN`

Bearer token for cluster auth

`KUBECFG_USERNAME`

Username for basic auth

`KUBECFG_PASSWORD`

Password for basic auth