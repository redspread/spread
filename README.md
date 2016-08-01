<p align="center"><img src="https://redspread.com/images/logo.svg" alt="logo" width= "400"/></p>

<p align="center"><a href="https://travis-ci.org/redspread/spread"><img alt="Build Status" src="https://travis-ci.org/redspread/spread.svg?branch=master" /></a> <a href="https://github.com/redspread/spread"><img alt="Release" src="https://img.shields.io/github/release/redspread/spread.svg" /></a> <a href="https://github.com/redspread/spread/blob/master/LICENSE"><img alt="Hex.pm" src="https://img.shields.io/hexpm/l/plug.svg" /></a> <a href="http://godoc.org/rsprd.com/spread"><img alt="GoDoc Status" src="https://godoc.org/rsprd.com/spread?status.svg" /></a></p>

<p align="center"><a href="https://redspread.com">Website</a> | <a href="http://redspread.readme.io">Docs</a> | <a href="http://slackin.redspread.com/">Slack</a> | <a href="mailto:founders@redspread.com">Email</a> | <a href="http://twitter.com/redspread">Twitter</a> | <a href="http://facebook.com/GetRedspread">Facebook</a></p>

#Spread: Git for Kubernetes

`spread` is a command line tool that makes it easy to version Kubernetes clusters, deploy to Kubernetes clusters in one command, and set up a local Kubernetes cluster (see: [localkube](https://github.com/redspread/localkube)). The project's goals are to:

* Guarantee reproducible Kubernetes deployments
* Be the fastest, simplest way to deploy Docker to production
* Enable collaborative deployment workflows that work well for one person or an entire team

See how we deployed Mattermost (<a href="https://github.com/redspread/kube-mattermost">and you can too!</a>):

<p align="center"><img src="http://i.imgur.com/Vohnd3e.gif" alt="logo" width= "800"/></p>

Spread is under open, active development. New features will be added regularly over the next few months - explore our [roadmap](./roadmap.md) to see what will be built next and send us pull requests for any features you’d like to see added.

See our [philosophy](./philosophy.md) for more on our mission and values. 

##Requirements
* Running Kubernetes cluster (<a href="https://blog.redspread.com/2016/02/04/google-container-engine-quickstart/">remote</a> or <a href="https://github.com/redspread/localkube">local</a>)
* [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
* [Go](https://golang.org/doc/install) (v 1.6)

##Installation

Install with `go get` (-d is for download only):

`go get -d rsprd.com/spread/cmd/spread`

Go into the correct directory:

`cd $GOPATH/src/rsprd.com/spread`

If libgit2 is not installed:

`make install-libgit2`

Then:

`make build/spread`

If an error about libraries missing comes up, set up your library path like:

`export LD_LIBRARY_PATH=/usr/local/lib:$ LD_LIBRARY_PATH`

Or, if you prefer using Homebrew (OS X only): 

`$ brew tap redspread/spread`  
`$ brew install spread-versioning`

##Git for Kubernetes

Spread versions your software environment (i.e. a Kubernetes cluster) like Git versions source code. Because Spread is built on top of libgit2, it takes advantage of Git's interface and functionality. This means after you deploy a Kubernetes object to a cluster, you can version the object by staging, commiting, and pushing it to a Spread repository. 

To get started, initialize Spread and set up a local Spread repository:

`spread init`

Here is our suggested workflow for versioning with Spread:

1. Create or edit your Kubernetes objects
2. Deploy your objects to a local or remote Kubernetes cluster. To use Spread's one-command deploy feature, make sure you've [set up your directory](https://github.com/redspread/spread/tree/versioning#faq) correctly, then `spread deploy .`
3. Stage an object: `spread add <objectType>/<objectName>`
4. Repeat until all objects have been staged
5. Commit your objects with a message: `spread commit -m "commit message"`
6. Check the status with `spread status` and diff with `spread diff`
7. Go ahead and try out the other commands - most Git commands, like `spread git log`, can be accessed using `spread git ...`

Spread versioning is highly experimental for the next few weeks. If you find any bugs or have any feature requests for Spread versioning, please file an issue, and know that the format for Spread may change! 

For more details on Spread commands, [see our docs](https://redspread.readme.io/docs/spread-commands).

##Spread Deploy Quickstart

Check out our <a href="https://redspread.readme.io/docs/getting-started">Getting Started Guide</a>.

##Localkube

Spread makes it easy to set up and iterate with [localkube](https://github.com/redspread/localkube), a local Kubernetes cluster streamlined for rapid development. 

**Requirements:**
* [Docker](https://docs.docker.com/engine/installation/)
* [docker-machine](https://docs.docker.com/machine/install-machine/)
* [VirtualBox](https://www.virtualbox.org/wiki/Downloads)

(Note: For Mac and Windows users, the fastest way to install everything above is [Docker Toolbox](https://www.docker.com/products/docker-toolbox).)

**Get started:**  

1. Create a machine called dev: `docker-machine create --driver virtualbox dev`
2. Start your docker-machine: `docker-machine start dev`
3. Connect to the docker daemon: `eval "$(docker-machine env dev)"`
4. Spin up a local cluster using [localkube](http://github.com/redspread/localkube): `spread cluster start`
5. To stop the cluster: `spread cluster stop`

**Suggested workflow:**
- `docker build` the image that you want to work with [1]
- Create Kubernetes objects that use the image build above
- Run `spread deploy .` to deploy to cluster [2]
- Iterate on your application, updating image and objects running `spread deploy .` each time you want to deploy changes
- To preview changes, grab the IP of your docker daemon with `docker-machine env <NAME>`and the returned `NodePort`, then put `IP:NodePort` in your browser
- When finished, run `spread cluster stop` to stop localkube
- To remove the container entirely, run `spread cluster stop -r`

[1] `spread` will soon integrate building ([#59](https://github.com/redspread/spread/issues/59))    
[2] Since `localkube` shares a Docker daemon with your host, there is no need to push images :)

[See more](https://github.com/redspread/localkube) for our suggestions when developing code with `localkube`.

##What's been done so far
 
* Spread versioning
* `spread deploy [-s] PATH [kubectl context]`: Deploys a Docker project to a Kubernetes cluster. It completes the following order of operations:
	* Reads context of directory and builds Kubernetes deployment hierarchy.
	* Updates all Kubernetes objects on a Kubernetes cluster.
	* Returns a public IP address, if type Load Balancer is specified. 
* [localkube](https://github.com/redspread/localkube): easy-to-setup local Kubernetes cluster for rapid development

##What's being worked on now

* Inner-app linking
* Parameterization
* [Redspread](redspread.com) (hosted Spread repository)

See more of our <a href="https://github.com/redspread/spread/blob/master/roadmap.md">roadmap</a> here!

##Future Goals
* Peer-to-peer syncing between local and remote Kubernetes clusters

##FAQ

**How are clusters selected?** Remote clusters are selected from the current kubectl context. Later, we will add functionality to explicitly state kubectl arguments. 

**How should I set up my directory?** In order to take advantage of Spread's one-command deploy feature, `spread deploy`, you'll need to set up your directory with a few specific naming conventions:

* All `ReplicationController` and `Pod` files should go in the root directory
* Any `ReplicationController` files should end in `.rc.yaml` or `.rc.json`, depending on the respective file extension
* Any `Pod` files should end in `.pod.yaml` or `.pod.json`, depending on the respective file extension
* All other Kubernetes object files should go in a directory named `rs`

There is no limit to the number of `ReplicationController`s or `Pod`s in the root directory.

Here is an example directory with Spread's naming conventions:

```Dockerfile
app.rc.yaml
database.rc.yaml
rs
 |_
    service.yaml
    secret.yaml
 ```
 
 **Why version objects instead of just files?** The object is the deterministic representation of state in Kubernetes. A useful analogy is "Kubernetes objects" are to "Docker images" like "Kubernetes object files" are to "Dockerfiles". By versioning the object itself, we can guarantee a 1:1 mapping with the Kubernetes cluster. This allows us to do things like diff two clusters and introduces future potential for linking between objects and repositories. 

##Contributing

We'd love to see your contributions - please see the CONTRIBUTING file for guidelines on how to contribute.

##Reporting bugs
If you haven't already, it's worth going through <a href="http://fantasai.inkedblade.net/style/talks/filing-good-bugs/">Elika Etemad's guide</a> for good bug reporting. In one sentence, good bug reports should be both *reproducible* and *specific*.

##Contact
Team: <a href="mailto:hello@redspread.com">hello@redspread.com</a>   
Slack: <a href="http://slackin.redspread.com">slackin.redspread.com</a>  
Planning: <a href="https://github.com/redspread/spread/blob/master/roadmap.md">Roadmap</a>  
Bugs: <a href="https://github.com/redspread/spread/issues">Issues</a>

##License
Spread is under the [Apache 2.0 license](https://tldrlegal.com/license/apache-license-2.0-(apache-2.0)). See the LICENSE file for details.
