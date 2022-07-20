package utils

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/transport"
)

// BuildClientsetFromConfig builds a Clientset from a given Kubeconfig.
//
// Takes an optional transport.WrapperFunc to add request/response middleware to
// api-server requests.
func BuildClientsetFromConfig(
	kubeconfig string,
	transportWrapper transport.WrapperFunc,
) (kubernetes.Interface, error) {
	restClientConfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cluster kubeconfig: %s", err)
	}

	if transportWrapper != nil {
		restClientConfig.Wrap(transportWrapper)
	}

	clientset, err := kubernetes.NewForConfig(restClientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build k8s client from Civo cluster kubeconfig: %s", err)
	}
	return clientset, nil
}
