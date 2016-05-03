###Spread roadmap
---

This is a high level roadmap for spread development. The dates below should not be considered authoritative, but rather indicative of a projected timeline for Spread. 

####Feb. 2016

* `spread deploy`: One command deploy to a remote Kubernetes cluster
	* Including building Docker context with deployment
	* Handle updates for all Kubernetes objects

####Mar. 2016

* Support for local development with Kubernetes
	* Easy setup of local cluster: see <a href="https://github.com/redspread/localkube">Localkube</a>
* Support for Linux and Windows

####April 2016

* Kubernetes configuration versioning
	* See our prototype: <a href="https://github.com/redspread/kit">Kit</a>

####May 2016

* Inner-app linking
* Easy rollbacks
* Local development
	* `spread build`: Builds Docker context and pushes to a local Kubernetes cluster.
	* `spread status`: Prints information about the state of the project

####June 2016

* Logging and debugging
	* `spread debug`: Returns the current spread version, Kubernetes version, Docker version, Go version, computer operating system, the Kubernetes config, spread config, and the content and timestamp of last command ran
	* `spread logs`: Returns logs for any deployment, automatic trying until logs are accessible.