![logo](../public/images/Redspread_Logo.png)

[![release](https://img.shields.io/badge/release-v0.0.1-red.svg)]() [![Hex.pm](https://img.shields.io/hexpm/l/plug.svg)]()



#spread: git for Docker deployment

###What is spread?

**spread is the fastest, simplest way to deploy Docker projects to Kubernetes clusters.** 

spread is a versioned container deployment workflow. It is under open, active development. New features will be added every two weeks over the next few months - explore our <a href="https://github.com/redspread/spread/blob/master/roadmap">roadmap</a> to see what will be built next and send us pull requests for any features you’d like to see added. 

The first feature is `spread deploy`, which enables users to deploy Docker to Kubernetes in one command. 

###Requirements
* <a href="https://docs.docker.com/engine/installation/">docker</a>
* <a href="https://docs.docker.com/machine/get-started/">docker-machine</a>
* <a href="https://blog.redspread.com/2016/02/04/google-container-engine-quickstart/">Kubernetes cluster with kubectl installed</a>

###Installation

`$ brew tap redspread/homebrew-spread`  
`$ brew install spread`

###Directory Structure

spread requires a specific project directory structure, as it builds from a hierarchy of entities:

* `Dockerfile`
* `*.ctr` - optional container file, there can be any number
* `pod.yaml` - pod file, there can be only one per directory
* `rc.yaml` - replication controller file, there can be only one per directory
* `/.k2e` - holds arbitrary Kubernetes objects, such as services and secrets

**Note:** The .ctr file is in .yaml format. Containers can still be placed in pods or replication controllers, but encouraging separate container files enables users to eventually reuse containers across an application.

###Hello World

This assumes you have a <a href="https://blog.redspread.com/2016/02/04/google-container-engine-quickstart/">running Kubernetes cluster</a> and <a href="https://docs.docker.com/machine/get-started/">docker-machine</a> installed.

1. Install spread with `$ brew tap redspread/homebrew-spread` then `$ brew install spread` 
2. Clone example project `$ git clone http://github.com/redspread/mvp-ex`
3. Start a docker machine `$ docker machine start <name>` or <a href="https://docs.docker.com/machine/get-started/">create a new machine</a>
4. Enter in your Docker registry configuration in the correct fields in .k2e/secret.yaml:
<pre><code>apiVersion: v1
kind: Secret
metadata:
  name: NAME
  namespace: NAMESPACE
data:
  .dockercfg: KEY
type: kubernetes.io/dockercfg</code></pre>
5. Build and deploy your project to Kubernetes: `$ spread deploy`
6. Grab the public IP and put it in your browser to see your website!

###Roadmap

* `spread status`: prints information about the state of Kubernetes objects and other entities 
* `spread debug`: returns info useful for debugging purposes
* `spread log`: returns logs for any pod, automatic trying until logs are accessible
* `spread build`: builds Docker context and pushes to a local Kubernetes cluster

See more of our <a href="https://github.com/redspread/spread/blob/master/roadmap">roadmap</a> here!

###Contributing

We'd love to see your contributions - please see CONTRIBUTING for guidelines on how to contribute.

###Contact
Founders: founders@redspread.com  
Slack: redspread.slack.com  
Planning/roadmap: <a href="http://github.com/redspread/spread/roadmap.md">roadmap</a>  
Bugs: <a href="https://github.com/redspread/spread/issues">issues</a>

###Reporting bugs
If you haven't already, it's worth going through <a href="http://fantasai.inkedblade.net/style/talks/filing-good-bugs/">Elika Etemad's guide</a> for good bug reporting.

###License
spread is under the Apache 2.0 license. See the LICENSE file for details.