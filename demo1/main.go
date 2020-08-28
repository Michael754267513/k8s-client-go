package main

import (
	"fmt"
	"github.com/owenliang/k8s-client-go/common"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getNamespacePodList(namespace string, clientset *kubernetes.Clientset) *core_v1.PodList {
	var podsList *core_v1.PodList
	var err error
	if podsList, err = clientset.CoreV1().Pods(namespace).List(meta_v1.ListOptions{}); err != nil {
		panic(err)
	}
	return podsList
}

func main() {
	var (
		clientset *kubernetes.Clientset
		//err error
	)

	// 初始化k8s客户端
	//if clientset, err = common.InitClient(); err != nil {
	//	goto FAIL
	//}
	clientset, _ = common.InitClient()
	pods := getNamespacePodList("kube-system", clientset)

	for _, v := range pods.Items {
		fmt.Printf("pods名称：%s 容器名称: %s 命名空间：%s  标签：%s 状态： %s  创建时间： %s  容器ip：%s\n", v.Name, v.Spec.Containers[0].Name, v.Namespace, v.Labels, v.Status.Phase, v.CreationTimestamp, v.Status.PodIP)
	}
	//// 获取default命名空间下的所有POD
	//if podsList, err = clientset.CoreV1().Pods("default").List(meta_v1.ListOptions{}); err != nil {
	//	goto FAIL
	//}
	////fmt.Println(podsList)
	////fmt.Println(podsList.Items)
	//for _,v := range  podsList.Items{
	//	fmt.Printf("pods名称：%s 容器名称: %s 命名空间：%s  标签：%s 状态： %s  创建时间： %s \n",v.Name,v.Spec.Containers[0].Name,v.Namespace,v.Labels,v.Status.Phase,v.CreationTimestamp)
	//}

	return

	//FAIL:
	//	fmt.Println(err)
	//	return
}
