package controller

import (
	"fmt"
	"github.com/owenliang/k8s-client-go/demo10/pkg/client/clientset/versioned"
	"github.com/owenliang/k8s-client-go/demo10/pkg/client/informers/externalversions/nginx_controller/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type NginxController struct {
	Clientset     *kubernetes.Clientset
	CrdClientset  *versioned.Clientset
	NginxInformer v1.NginxInformer

	NginxWorkqueue workqueue.RateLimitingInterface
}

func (nginxController *NginxController) Start() (err error) {
	var (
		stopCh = make(chan struct{})
		i      int
		syncOk bool
	)

	// nginx informer的event handler
	nginxController.NginxInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			nginxController.OnAddNginx(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			nginxController.OnUpdateNginx(oldObj, newObj)
		},
		DeleteFunc: func(obj interface{}) {
			nginxController.OnDeleteNginx(obj)
		},
	})

	//  event handler会把event丢到workqueue里, 被processor消费
	nginxController.NginxWorkqueue = workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Nginx")

	// nginx informer开始拉取事件，存到local cache，并回调event handler
	go nginxController.NginxInformer.Informer().Run(stopCh)

	// 等待etcd已有数据都下载回来, 再启动事件处理线程, 这样local cache可以反馈出贴近准实时的etcd数据，供逻辑决策准确
	if syncOk = cache.WaitForCacheSync(stopCh, nginxController.NginxInformer.Informer().HasSynced); !syncOk {
		err = fmt.Errorf("sync失败")
		return
	}

	// 启动nginx event processor
	for i = 0; i < 2; i++ {
		go nginxController.runNginxWorker()
	}

	return
}
