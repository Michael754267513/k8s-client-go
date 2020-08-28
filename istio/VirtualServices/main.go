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
	istioClient.NetworkingV1alpha3().VirtualServices(namespace).Delete(context.TODO(), "vs-test", v1.DeleteOptions{})

	var (
		httpRouteList            []*networkingv1alpha3.HTTPRoute
		HTTPRouteDestinationList []*networkingv1alpha3.HTTPRouteDestination
	)

	HTTPRouteDestination := &networkingv1alpha3.HTTPRouteDestination{
		Destination: &networkingv1alpha3.Destination{
			Host:   "reviews",
			Subset: "v2",
		},
		Weight: 50,
	}
	HTTPRouteDestinationList = append(HTTPRouteDestinationList, HTTPRouteDestination)
	httpRouteSign := networkingv1alpha3.HTTPRoute{

		Route: HTTPRouteDestinationList,
	}
	httpRouteList = append(httpRouteList, &httpRouteSign)
	virtualService := &v1alpha3.VirtualService{
		ObjectMeta: v1.ObjectMeta{
			Name:      "vs-test",
			Namespace: namespace,
		},
		Spec: networkingv1alpha3.VirtualService{
			Hosts: []string{"reviews"},
			//Gateways: []string{""},
			Http: httpRouteList,
		},
	}

	vs, err := istioClient.NetworkingV1alpha3().VirtualServices(namespace).Create(context.TODO(), virtualService, v1.CreateOptions{})
	if err != nil {
		return
	}
	log.Print(vs)

	//vsList,err := istioClient.NetworkingV1alpha3().VirtualServices(namespace).List(context.TODO(),v1.ListOptions{})
	//if err != nil {
	//	return
	//}
	//for i := range vsList.Items {
	//	vs := vsList.Items[i]
	//
	//	log.Printf("VirtualService Hosts: %+v -- VirtualService Gateway:  %+v -- VirtualService http:   %+v",vs.Spec.Hosts,vs.Spec.Gateways,vs.Spec.Http)
	//
	//}
}
