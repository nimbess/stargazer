package client

import (
	"github.com/nimbess/stargazer/pkg/crd/api/unp/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func (c *UnpConfigV1Alpha1Client) UnpConfigs(namespace string) UnpConfigInterface {
	return &unpConfigclient{
		client: c.restClient,
		ns:     namespace,
	}
}

type UnpConfigV1Alpha1Client struct {
	restClient rest.Interface
}

type UnpConfigInterface interface {
	Create(obj *v1.UnpConfig) (*v1.UnpConfig, error)
	Update(obj *v1.UnpConfig) (*v1.UnpConfig, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	Get(name string) (*v1.UnpConfig, error)
}

type unpConfigclient struct {
	client rest.Interface
	ns     string
}

func (c *unpConfigclient) Create(obj *v1.UnpConfig) (*v1.UnpConfig, error) {
	result := &v1.UnpConfig{}
	err := c.client.Post().
		Namespace(c.ns).Resource("unpconfigs").
		Body(obj).Do().Into(result)
	return result, err
}

func (c *unpConfigclient) Update(obj *v1.UnpConfig) (*v1.UnpConfig, error) {
	result := &v1.UnpConfig{}
	err := c.client.Put().
		Namespace(c.ns).Resource("unpconfigs").
		Body(obj).Do().Into(result)
	return result, err
}

func (c *unpConfigclient) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).Resource("unpconfigs").
		Name(name).Body(options).Do().
		Error()
}

func (c *unpConfigclient) Get(name string) (*v1.UnpConfig, error) {
	result := &v1.UnpConfig{}
	err := c.client.Get().
		Namespace(c.ns).Resource("unpconfigs").
		Name(name).Do().Into(result)
	return result, err
}
