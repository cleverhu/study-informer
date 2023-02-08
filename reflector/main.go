package main

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"

	"github.com/cleverhu/study-informer/lib"
)

func main() {

	cs := lib.InitClient()
	df := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{KeyFunction: cache.MetaNamespaceKeyFunc})
	lwClient := cache.NewListWatchFromClient(cs.CoreV1().RESTClient(), "pods", "default", fields.Everything())
	reflector := cache.NewReflector(lwClient, &corev1.Pod{}, df, 0)
	go func() {
		reflector.Run(wait.NeverStop)
	}()

	for {
		obj, _ := df.Pop(func(i interface{}) error {
			for _, delta := range i.(cache.Deltas) {
				pod := delta.Object.(*corev1.Pod)
				fmt.Println(delta.Type, pod.Name, pod.Namespace, pod.Status.Phase)
			}
			return nil
		})
		fmt.Printf("%T\n", obj)
	}
}
