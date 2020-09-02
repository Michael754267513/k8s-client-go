package main

import (
	//"net/http"
	//
	//"github.com/emicklei/go-restful"
	//restfulspec "github.com/emicklei/go-restful-openapi"

	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "app.k8s.io"

var GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1beta1"}

//var (
//WebServiceBuilder = runtime.NewContainerBuilder(addWebService)
//AddToContainer    = WebServiceBuilder.AddToContainer
//)

//func addWebService(c *restful.Container) error {
//
//ok := "ok"
//webservice := runtime.NewWebService(GroupVersion)

func main() {
	fmt.Println(GroupVersion)
}
