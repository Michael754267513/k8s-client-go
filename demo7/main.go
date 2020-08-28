package main

import (
	"fmt"
	"github.com/owenliang/k8s-client-go/common"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	var (
		clientset *kubernetes.Clientset
		tailLines int64
		req       *rest.Request
		res       rest.Result
		logs      []byte
		err       error
	)

	// 初始化k8s客户端
	if clientset, err = common.InitClient(); err != nil {
		goto FAIL
	}

	// 生成获取POD日志请求
	req = clientset.CoreV1().Pods("kube-system").GetLogs("kubernetes-dashboard-5c7687cf8-5wkrm", &core_v1.PodLogOptions{Container: "kubernetes-dashboard", TailLines: &tailLines})

	// req.Stream()也可以实现Do的效果

	// 发送请求
	if res = req.Do(); res.Error() != nil {
		err = res.Error()
		goto FAIL
	}

	// 获取结果
	if logs, err = res.Raw(); err != nil {
		goto FAIL
	}

	fmt.Println("容器输出:", string(logs))
	return

FAIL:
	fmt.Println(err)
	return
}
