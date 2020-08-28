package main

import (
	"context"
	"log"

	"github.com/owenliang/k8s-client-go/common"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	namespace := "bookinfo"
	restConfig, err := common.GetRestConf()
	if err != nil {
		return
	}
	istioClient, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		return
	}
	gw, err := istioClient.NetworkingV1alpha3().Gateways(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get Gateway in %s namespace: %s", namespace, err)
	}

	for i := range gw.Items {
		gw := gw.Items[i]
		for _, s := range gw.Spec.GetServers() {
			log.Printf("Index: %d Gateway servers: %+v\n", i, s)
		}
	}

	var (
		gateway  *v1alpha3.Gateway
		services []*networkingv1alpha3.Server
	)

	service := &networkingv1alpha3.Server{
		Port: &networkingv1alpha3.Port{
			Number: 80,
			// MUST BE one of HTTP|HTTPS|GRPC|HTTP2|MONGO|TCP|TLS.
			Protocol: "HTTP",
			Name:     "http",
		},
		// $hide_from_docs
		// The ip or the Unix domain socket to which the listener should be bound
		// to. Format: `x.x.x.x` or `unix:///path/to/uds` or `unix://@foobar`
		// (Linux abstract namespace). When using Unix domain sockets, the port
		// number should be 0.
		//Bind:                 "",
		Hosts:                []string{"*.hzeng.com.cn"},
		Tls:                  nil,
		DefaultEndpoint:      "",
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}

	services = append(services, service)
	gateway = &v1alpha3.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name:      "gw-bookinfo",
			Namespace: namespace,
		},
		Spec: networkingv1alpha3.Gateway{
			Servers: services,
			//Selector:             nil,

		},
	}

	istioClient.NetworkingV1alpha3().Gateways(namespace).Create(context.TODO(), gateway, v1.CreateOptions{})

	istioClient.NetworkingV1alpha3().Gateways(namespace).Delete(context.TODO(), "gw-bookinfo", v1.DeleteOptions{})
}
