package pods

import (
	"fmt"

	"github.com/owenliang/k8s-client-go/common"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Pods struct {
	ContainerName string
	NameSpace     string
	PodName       string
	PodiP         string
	status        string
}

func ListNamespacesPods(namespace string) {
	var (
		clientset *kubernetes.Clientset
		podsList  *core_v1.PodList
		//resPodList []Pods
		pods Pods
		err  error
	)
	if clientset, err = common.InitClient(); err != nil {
		goto ERROR
	}
	if podsList, err = clientset.CoreV1().Pods(namespace).List(meta_v1.ListOptions{}); err != nil {
		goto ERROR
	}
	for _, v := range podsList.Items {
		fmt.Printf("pods名称：%s 容器名称: %s 命名空间：%s  标签：%s 状态： %s  创建时间： %s  容器ip：%s\n", v.Name, v.Spec.Containers[0].Name, v.Namespace, v.Labels, v.Status.Phase, v.CreationTimestamp, v.Status.PodIP)
		pods.ContainerName = v.Spec.Containers[0].Name
		pods.NameSpace = v.Namespace
		pods.PodName = v.Name
		pods.PodiP = v.Status.PodIP
		pods.status = string(v.Status.Phase)
		//resPodList = append(resPodList,pods)
		fmt.Println(pods)
	}
	//return resPodList
ERROR:
	panic(err)
}
