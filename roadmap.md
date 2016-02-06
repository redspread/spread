###spread roadmap
---

This is a high level roadmap for spread development. The dates below should not be considered authoritative, but rather indicative of a projected timeline for spread. 

####Feb. 2016

* Support deploying multiple containers to Kubernetes
	* `spread deploy`: One command to deploy to a remote cluster
	* Entity linking
* Logging and debugging
	* `spread debug`: Returns the current spread version, Kubernetes version, Docker version, Go version, computer operating system, the Kubernetes config, spread config, and the content and timestamp of last command ran
	* `spread logs`: Prints logs from deployment of current project

####Mar. 2016

* Support for local development with Kubernetes
	* Easy setup of local cluster
	* Cluster-to-cluster syncing
	* `spread build`: One command to build to a local cluster
	* `spread status`: Prints information about the state of the project