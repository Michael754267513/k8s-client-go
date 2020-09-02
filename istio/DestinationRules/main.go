package main

import (
	"context"
	"log"

	"github.com/gogo/protobuf/types"
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
	istioClient.NetworkingV1alpha3().DestinationRules(namespace).Delete(context.TODO(), "dr-reviews", v1.DeleteOptions{})
	dr, err := istioClient.NetworkingV1alpha3().DestinationRules(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return
	}
	for i := range dr.Items {
		dr := dr.Items[i]
		log.Printf("DestinationRule : %+v\n", dr.Spec)
	}
	var (
		destinationRule *v1alpha3.DestinationRule
		subsetList      []*networkingv1alpha3.Subset
	)

	// 设置subset
	subset := &networkingv1alpha3.Subset{
		Name:   "v2",
		Labels: map[string]string{"version": "v2"},
		//TrafficPolicy:        nil,

	}
	subsetList = append(subsetList, subset)

	destinationRule = &v1alpha3.DestinationRule{
		TypeMeta: v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Namespace: namespace,
			Name:      "dr-reviews",
		},
		Spec: networkingv1alpha3.DestinationRule{
			Host:    "reviews",
			Subsets: subsetList,
			TrafficPolicy: &networkingv1alpha3.TrafficPolicy{
				LoadBalancer: &networkingv1alpha3.LoadBalancerSettings{
					LbPolicy: &networkingv1alpha3.LoadBalancerSettings_Simple{
						Simple: networkingv1alpha3.LoadBalancerSettings_PASSTHROUGH,
					},
					LocalityLbSetting: nil,
				},
				ConnectionPool: &networkingv1alpha3.ConnectionPoolSettings{
					Tcp: &networkingv1alpha3.ConnectionPoolSettings_TCPSettings{
						// Maximum number of HTTP1 /TCP connections to a destination host. Default 2^32-1.
						MaxConnections: 200,
						// TCP connection timeout. format: 1h/1m/1s/1ms. MUST BE >=1ms. Default is 10s.
						ConnectTimeout: nil,
						//If set then set SO_KEEPALIVE on the socket to enable TCP Keepalives.
						TcpKeepalive: &networkingv1alpha3.ConnectionPoolSettings_TCPSettings_TcpKeepalive{
							/*
									net.ipv4.tcp_keepalive_intvl = 75
									net.ipv4.tcp_keepalive_probes = 9
									net.ipv4.tcp_keepalive_time = 7200
								   对应下列参数
							*/
							// 默认探测次数 默认是9(通道没有请求的时候进行探测)
							Probes: 9,
							// 默认是7200s  操作系统kernel
							Time: &types.Duration{
								Seconds: 7200,
								Nanos:   0,
							},
							Interval: &types.Duration{
								Seconds: 75,
								Nanos:   0,
							},
						},
					},
					Http: &networkingv1alpha3.ConnectionPoolSettings_HTTPSettings{
						// Maximum number of pending HTTP requests to a destination. Default 2^32-1.
						// 最大请求数
						Http1MaxPendingRequests: 200,
						// Maximum number of requests to a backend. Default 2^32-1.
						// 每个后端最大请求数
						Http2MaxRequests: 20,
						// Maximum number of requests per connection to a backend. Setting this
						// parameter to 1 disables keep alive. Default 0, meaning "unlimited",
						// up to 2^29.
						// 是否启用keepalive对后端进行长链接 0 表示启用
						MaxRequestsPerConnection: 0,
						// Maximum number of retries that can be outstanding to all hosts in a
						// cluster at a given time. Defaults to 2^32-1.
						// 在给定时间内最大的重试次数
						MaxRetries: 1,
						// The idle timeout for upstream connection pool connections. The idle timeout is defined as the period in which there are no active requests.
						// If not set, the default is 1 hour. When the idle timeout is reached the connection will be closed.
						// Note that request based timeouts mean that HTTP/2 PINGs will not keep the connection alive. Applies to both HTTP1.1 and HTTP2 connections.
						// 不设置默认1小时没有请求，断开后端连接
						IdleTimeout: nil,
						//// Specify if http1.1 connection should be upgraded to http2 for the associated destination.
						//H2UpgradePolicy:          0,
					},
				},
				// 类似nginx的 next upstream
				//OutlierDetection:     nil,
				//Tls:                  nil,
				PortLevelSettings: nil,
			},
		},
	}
	dr1, err := istioClient.NetworkingV1alpha3().DestinationRules(namespace).Create(context.TODO(), destinationRule, v1.CreateOptions{})
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(dr1)
}
