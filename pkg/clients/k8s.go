package k8s

import (
	"context"

	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient interface {
	GetSecret(ctx context.Context, namespace, name string) (*core_v1.Secret, error)
}

type k8sClient struct {
	clientSet *kubernetes.Clientset
}

func (k *k8sClient) GetSecret(ctx context.Context, namespace, name string) (*core_v1.Secret, error) {
	return k.clientSet.CoreV1().Secrets(namespace).Get(ctx, name, meta_v1.GetOptions{})
}

func NewK8sClient(kcfgPath string) (K8sClient, error) {
	rCfg, err := clientcmd.BuildConfigFromFlags("", kcfgPath)
	if err != nil {
		return nil, err
	}

	c, err := kubernetes.NewForConfig(rCfg)
	if err != nil {
		return nil, err
	}

	return &k8sClient{clientSet: c}, nil
}
