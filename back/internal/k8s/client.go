package k8s

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Client interface {
	ListResources(ctx context.Context, apiVersion, kind, namespace string, selector map[string]string) (*unstructured.UnstructuredList, error)
}

type client struct {
	dynamic   dynamic.Interface
	discovery discovery.DiscoveryInterface
	mapper    *restmapper.DeferredDiscoveryRESTMapper
	logger    *slog.Logger
}

func NewClient(logger *slog.Logger) (Client, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	dyClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	discClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create discovery client: %w", err)
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discClient))

	return &client{
		dynamic:   dyClient,
		discovery: discClient,
		mapper:    mapper,
		logger:    logger,
	}, nil
}

func (c *client) ListResources(ctx context.Context, apiVersion, kind, namespace string, selector map[string]string) (*unstructured.UnstructuredList, error) {
	// Parse GroupVersion
	gv, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse apiVersion '%s': %w", apiVersion, err)
	}

	// Find GVR
	gvk := gv.WithKind(kind)
	mapping, err := c.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to find REST mapping for %s: %w", gvk, err)
	}

	var res dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace && namespace != "" {
		res = c.dynamic.Resource(mapping.Resource).Namespace(namespace)
	} else {
		res = c.dynamic.Resource(mapping.Resource)
	}

	opts := metav1.ListOptions{}
	if len(selector) > 0 {
		ls := metav1.LabelSelector{MatchLabels: selector}
		labelSelector, err := metav1.LabelSelectorAsSelector(&ls)
		if err != nil {
			return nil, fmt.Errorf("failed to convert label selector: %w", err)
		}
		opts.LabelSelector = labelSelector.String()
	}

	list, err := res.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}

	return list, nil
}
