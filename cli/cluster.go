package cli

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/mitchellh/go-homedir"
	kubectlapi "k8s.io/kubernetes/pkg/client/unversioned/clientcmd/api"
	kubectlcfg "k8s.io/kubernetes/pkg/kubectl/cmd/config"
)

const (
	LocalkubeContainerName = "/localkube"
	LocalkubeImageName     = "redspreadapps/localkube"
	LocalkubeDefaultTag    = "latest"

	DefaultHostDataDir = "~/.localkube/data"
	ContainerDataDir   = "/var/localkube/data"
	KubectlName        = "localkube"
)

// Cluster manages the localkube Kubernetes development environment.
func (s SpreadCli) Cluster() *cli.Command {
	return &cli.Command{
		Name:        "cluster",
		Usage:       "spread cluster [-a] [-t <tag>] <start|stop> [ClusterDataDirectory]",
		Description: "Manages localkube Kubernetes development environment",
		ArgsUsage:   "-a will attach to the process and print logs to stdout, -t specfies localkube tag to use, default is latest.",
		Action: func(c *cli.Context) {
			action := strings.ToLower(c.Args().First())
			switch {
			case "start" == action:
				s.startLocalkube(c)
			case "stop" == action:
				s.stopLocalkube(c)
			default:
				s.printf("Invalid option `%s`, must choose start or stop", action)
			}
		},
	}
}

func (s SpreadCli) startLocalkube(c *cli.Context) {
	client := s.dockerOrErr()

	dataDir := c.Args().Get(1)
	if len(dataDir) == 0 {
		var err error
		dataDir, err = homedir.Expand(DefaultHostDataDir)
		if err != nil {
			s.fatalf("Unable to expand home directory: %v", err)
		}
	}

	tag := c.String("t")
	if len(tag) == 0 {
		tag = LocalkubeDefaultTag
	}

	ctrOpts := localkube(c.Bool("a"), dataDir, tag)
	ctr, err := client.CreateContainer(ctrOpts)
	if err != nil {
		if err.Error() == "no such image" {
			s.printf("Pulling localkube image...")
			err = client.PullImage(docker.PullImageOptions{
				Repository: LocalkubeImageName,
				Tag:        tag,
			}, docker.AuthConfiguration{})
			if err != nil {
				s.fatalf("Failed to pull localkube image: %v", err)
			}

			s.startLocalkube(c)
			return
		} else if err.Error() == "container already exists" {
			// replace container if already exists
			err = client.RemoveContainer(docker.RemoveContainerOptions{
				ID: LocalkubeContainerName,
			})
			if err != nil {
				s.fatalf("Failed to start container: %v", err)
			}
			s.startLocalkube(c)
		}
		s.fatalf("Failed to create localkube container: %v", err)

	}

	binds := []string{
		"/sys:/sys:rw",
		"/var/lib/docker:/var/lib/docker",
		"/mnt/sda1/var/lib/docker:/mnt/sda1/var/lib/docker",
		"/var/lib/kubelet:/var/lib/kubelet",
		"/var/run:/var/run:rw",
		"/:/rootfs:ro",
	}

	// if provided mount etcd data dir
	if len(dataDir) != 0 {
		dataBind := fmt.Sprintf("%s:%s", dataDir, ContainerDataDir)
		binds = append(binds, dataBind)
	}

	hostConfig := &docker.HostConfig{
		Binds:         binds,
		NetworkMode:   "host",
		RestartPolicy: docker.AlwaysRestart(),
		PidMode:       "host",
		Privileged:    true,
	}
	err = client.StartContainer(ctr.ID, hostConfig)
	if err != nil {
		s.fatalf("Failed to start localkube: %v", err)
	}

	s.printf("Started localkube...")

	err = setupContext(client.Endpoint())
	if err != nil {
		s.fatalf("Could not configure kubectl context.")
	}
	s.printf("Setup and using kubectl `%s` context.", KubectlName)
	return
}

func (s SpreadCli) stopLocalkube(c *cli.Context) {
	client := s.dockerOrErr()

	ctrs, err := client.ListContainers(docker.ListContainersOptions{
		All: true,
		Filters: map[string][]string{
			"label": []string{"rsprd.com/name=localkube"},
		},
	})
	if err != nil {
		s.fatalf("Could not list containers: %v", err)
	}

	for _, ctr := range ctrs {
		if strings.HasPrefix(ctr.Status, "Up") {
			s.printf("Stopping container `%s`...\n", ctr.ID)
			if err := client.StopContainer(ctr.ID, 5); err != nil {
				s.fatalf("Could not kill container: %v", err)
			}
		}

		if err := client.RemoveContainer(docker.RemoveContainerOptions{ID: ctr.ID}); err != nil {
			s.fatalf("Could not remove container: %v", err)
		}
	}
}

func (s SpreadCli) dockerOrErr() *docker.Client {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		s.fatalf("Could not create Docker client: %v", err)
	}

	_, err = client.Version()
	if err != nil {
		s.fatalf("Unable to establish connection with Docker daemon: %v", err)
	}
	return client
}

func localkube(attach bool, dataDir, tag string) docker.CreateContainerOptions {
	return docker.CreateContainerOptions{
		Name: LocalkubeContainerName,
		Config: &docker.Config{
			Hostname:     "localkube",
			AttachStderr: attach,
			AttachStdout: attach,
			Image:        fmt.Sprintf("%s:%s", LocalkubeImageName, tag),
			Env: []string{
				fmt.Sprintf("KUBE_ETCD_DATA_DIRECTORY=%s", ContainerDataDir),
			},
			Labels: map[string]string{
				"rsprd.com/name": "localkube",
			},
			StopSignal: "SIGINT",
		},
	}
}

func identifyHost(endpoint string) (string, error) {
	beginPort := strings.LastIndex(endpoint, ":")
	switch {
	// if using TCP use provided host
	case strings.HasPrefix(endpoint, "tcp://"):
		return endpoint[6:beginPort], nil
	// assuming localhost if Unix
	// TODO: Make this customizable
	case strings.HasPrefix(endpoint, "unix://"):
		return "127.0.0.1", nil
	}
	return "", fmt.Errorf("Could not determine localkube API server from endpoint `%s`", endpoint)
}

func setupContext(endpoint string) error {
	host, err := identifyHost(endpoint)
	if err != nil {
		return fmt.Errorf("Could not identify host: %v", err)
	}

	pathOpts := kubectlcfg.NewDefaultPathOptions()

	config, err := pathOpts.GetStartingConfig()
	if err != nil {
		return fmt.Errorf("could not setup config: %v", err)
	}

	cluster, exists := config.Clusters[KubectlName]
	if !exists {
		cluster = kubectlapi.NewCluster()
	}

	// configure cluster
	cluster.Server = fmt.Sprintf("%s:8080", host)
	cluster.InsecureSkipTLSVerify = true
	config.Clusters[KubectlName] = cluster

	context, exists := config.Contexts[KubectlName]
	if !exists {
		context = kubectlapi.NewContext()
	}

	// configure context
	context.Cluster = KubectlName
	config.Contexts[KubectlName] = context

	config.CurrentContext = KubectlName

	return kubectlcfg.ModifyConfig(pathOpts, *config, true)
}
