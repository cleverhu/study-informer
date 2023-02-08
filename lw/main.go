package main

import (
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"

	"github.com/cleverhu/study-informer/lib"
)

func main() {
	cs := lib.InitClient()

	lwClient := cache.NewListWatchFromClient(cs.CoreV1().RESTClient(), "pods", "default", fields.Everything())
	list, err := lwClient.List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#T\n", list)
	lw, err := lwClient.Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case data, ok := <-lw.ResultChan():
			if ok {
				fmt.Printf("Obj: %#v, Type: %#v\n", data.Object.(*corev1.Pod).Name, data.Type)
			}
		}
	}
}
