<center><img src="https://redspread.com/images/logo.svg" alt="logo" width= "400"/>

[![Build Status](https://travis-ci.org/redspread/spread.svg?branch=master)](https://travis-ci.org/redspread/spread) [![release](https://img.shields.io/badge/release-v0.0.3-red.svg)]() [![Hex.pm](https://img.shields.io/hexpm/l/plug.svg)]() [![](https://godoc.org/rsprd.com/spread?status.svg)](http://godoc.org/rsprd.com/spread)</center>

<center>[Website](https://redspread.com) | [Slack](http://redspread.slack.com) | <a href="mailto:founders@redspread.com">Email</a> | <a href="http://twitter.com/redspread">Twitter</a> | <a href="http://facebook.com/GetRedspread">Facebook</a></center>

#Docker to Kubernetes in one command

`spread` is a command line tool that builds and deploys a [Docker](https://github.com/docker/docker) project to a [Kubernetes](https://github.com/kubernetes/kubernetes) cluster in one command. The project's goals are to:

* Enable rapid iteration with Kubernetes
* Be the fastest, simplest way to deploy Docker to production
* Work well for a single developer or an entire team (no more broken bash scripts!)


Spread is under open, active development. New features will be added regularly over the next few months - explore our [roadmap](./roadmap.md) to see what will be built next and send us pull requests for any features you’d like to see added.

See our [philosophy](./philosophy.md) for more on our mission and values.

## Requirements
* [Kubernetes cluster with kubectl installed](https://blog.redspread.com/2016/02/04/google-container-engine-quickstart/)

## Installation

For OSX
`$ brew tap redspread/spread`
`$ brew install spread`


## Quickstart

This assumes you have a [running Kubernetes cluster](https://blog.redspread.com/2016/02/04/google-container-engine-quickstart/).

1. Install Spread with `$ brew tap redspread/spread` then `$ brew install spread`
2. Clone [Mattermost](http://mattermost.com), the open source Slack `$ git clone http://github.com/redspread/kube-mattermost` and change into the newly created directory `$ cd mattermost`.
5. Deploy Mattermost to Kubernetes: `$ spread deploy .`
6. Grab the public IP from the output and open it in your browser to see your self-hosted app!

For a more detailed walkthrough, see the full [guide](https://github.com/redspread/kube-mattermost).

## What's been done so far

* `$ spread deploy [-s] PATH [kubectl context]`: Deploys a Docker project to a Kubernetes cluster. It completes the following order of operations:
	* Reads context of directory and builds Kubernetes deployment hierarchy.
	* Updates all Kubernetes objects on a Kubernetes cluster.
	* Returns a public IP address, if type Load Balancer is specified.
* Established an implicit hierarchy of Kubernetes objects
* Multi-container deployment

## What's being worked on now

* Build functionality for `$ spread deploy` so it also builds any images indicated to be built and pushes those images to the indicated Docker registry.
* `$ spread deploy -p`: Pushes all images to registry, even those not built by `$ spread deploy`.
* Support for Linux and Windows
* Inner-app linking
* `$ spread logs`: Returns logs for any deployment, automatic trying until logs are accessible.
* `$ spread build`: Builds Docker context and pushes to a local Kubernetes cluster.
* `$ spread rewind`: Quickly rollback to a previous deployment.

See more of our [roadmap here!](https://github.com/redspread/spread/blob/master/roadmap.md)

## Future Goals

* Develop workflow for container versioning (containers = image + config)
* Introduce paramaterization for container configuration

## FAQ

**How are clusters selected?** Remote clusters are selected from the current kubectl context. Later, we will add functionality to explicitly state kubectl arguments.

**How should I set up my directory?** Spread requires a specific project directory structure, as it builds from a hierarchy of entities:

* `Dockerfile`
* `*.ctr` - optional container file, there can be any number
* `pod.yaml` - pod file, there can be only one per directory
* `rc.yaml` - replication controller file, there can be only one per directory
* `/.k2e` - holds arbitrary Kubernetes objects, such as services and secrets

**What is the *.ctr file?** The .ctr file is the container struct usually found in the pod.yaml or rc.yaml. Containers can still be placed in pods or replication controllers, but we're encouraging separate container files because it enables users to eventually reuse containers across an application.

**Can I deploy a project with just a Dockerfile and *.ctr?** Yes. Spread implicitly infers the rest of the app hierarchy.

## Contributing

We'd love to see your contributions - please see the CONTRIBUTING file for guidelines on how to contribute.

## Reporting bugs
Please use our [issue template](https://github.com/redspread/spread/blob/master/ISSUE_TEMPLATE.md).

## Contact
Founders: [founders@redspread.com](mailto:founders@redspread.com)

Slack: [redspread.slack.com](http://redspread.slack.com)

Planning/roadmap: [roadmap](http://github.com/redspread/spread/roadmap.md)

Bugs: [issues](https://github.com/redspread/spread/issues)

## License
[Spread is under the Apache 2.0 license.](https://github.com/redspread/spread/blob/master/LICENSE)
