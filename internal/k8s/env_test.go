package k8s

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd/api"
)

type expectedConfig struct {
	expected Env
	input    *api.Config
}

func TestProvideEnv(t *testing.T) {
	minikubeContexts := map[string]*api.Context{
		"minikube": &api.Context{
			Cluster: "minikube",
		},
	}
	dockerDesktopContexts := map[string]*api.Context{
		"docker-for-desktop": &api.Context{
			Cluster: "docker-for-desktop-cluster",
		},
	}
	dockerDesktopEdgeContexts := map[string]*api.Context{
		"docker-for-desktop": &api.Context{
			Cluster: "docker-desktop",
		},
	}
	gkeContexts := map[string]*api.Context{
		"gke_blorg-dev_us-central1-b_blorg": &api.Context{
			Cluster: "gke_blorg-dev_us-central1-b_blorg",
		},
	}
	kind5Contexts := map[string]*api.Context{
		"kubernetes-admin@kind-1": &api.Context{
			Cluster: "kind",
		},
	}
	microK8sContexts := map[string]*api.Context{
		"microk8s": &api.Context{
			Cluster: "microk8s-cluster",
		},
	}

	homedir, err := homedir.Dir()
	assert.NoError(t, err)
	k3dContexts := map[string]*api.Context{
		"default": &api.Context{
			LocationOfOrigin: filepath.Join(homedir, ".config", "k3d", "k3s-default", "kubeconfig.yaml"),
			Cluster:          "default",
		},
	}
	kind5NamedClusterContexts := map[string]*api.Context{
		"default": &api.Context{
			LocationOfOrigin: filepath.Join(homedir, ".kube", "kind-config-integration"),
			Cluster:          "integration",
		},
	}
	kind6Contexts := map[string]*api.Context{
		"kind-custom-name": &api.Context{
			Cluster: "kind-custom-name",
		},
	}
	table := []expectedConfig{
		{EnvNone, &api.Config{}},
		{EnvUnknown, &api.Config{CurrentContext: "aws"}},
		{EnvMinikube, &api.Config{CurrentContext: "minikube", Contexts: minikubeContexts}},
		{EnvDockerDesktop, &api.Config{CurrentContext: "docker-for-desktop", Contexts: dockerDesktopContexts}},
		{EnvDockerDesktop, &api.Config{CurrentContext: "docker-for-desktop", Contexts: dockerDesktopEdgeContexts}},
		{EnvGKE, &api.Config{CurrentContext: "gke_blorg-dev_us-central1-b_blorg", Contexts: gkeContexts}},
		{EnvKIND5, &api.Config{CurrentContext: "kubernetes-admin@kind-1", Contexts: kind5Contexts}},
		{EnvMicroK8s, &api.Config{CurrentContext: "microk8s", Contexts: microK8sContexts}},
		{EnvK3D, &api.Config{CurrentContext: "default", Contexts: k3dContexts}},
		{EnvKIND5, &api.Config{CurrentContext: "default", Contexts: kind5NamedClusterContexts}},
		{EnvKIND6, &api.Config{CurrentContext: "kind-custom-name", Contexts: kind6Contexts}},
	}

	for _, tt := range table {
		t.Run(tt.input.CurrentContext, func(t *testing.T) {
			actual := ProvideEnv(context.Background(), tt.input)
			if actual != tt.expected {
				t.Errorf("Expected %s, actual %s", tt.expected, actual)
			}
		})
	}
}
