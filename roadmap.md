###spread roadmap
---

This is a high level roadmap for spread development. The dates below should not be considered authoritative, but rather indicative of a projected timeline for spread. 

####Feb. 2016

* `spread deploy`: One command deploy to a remote Kubernetes cluster
	* Including building Docker context with deployment
	* `spread deploy -p`: Pushes all images to registry, even those not built by `spread deploy`.
	* Handle updates for all Kubernetes objects

####Mar. 2016

* Logging and debugging
	* `spread debug`: Returns the current spread version, Kubernetes version, Docker version, Go version, computer operating system, the Kubernetes config, spread config, and the content and timestamp of last command ran
	* `spread logs`: Returns logs for any deployment, automatic trying until logs are accessible.
* Support for local development with Kubernetes
	* Easy setup of local cluster
	* Cluster-to-cluster syncing
	* `spread build`: Builds Docker context and pushes to a local Kubernetes cluster.
	* `spread status`: Prints information about the state of the project
* Inner-app linking
* Support for Linux and Windows

####Apr. 2016

* Container versioning
	* Version the image + configuration of containers 
	* `spread rewind`: Quickly rollback to a previous deployment.