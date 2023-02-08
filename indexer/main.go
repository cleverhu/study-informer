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
	store := cache.NewStore(cache.MetaNamespaceKeyFunc)
	df := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{
		KeyFunction:  cache.MetaNamespaceKeyFunc,
		KnownObjects: store,
	})
	lwClient := cache.NewListWatchFromClient(cs.CoreV1().RESTClient(), "pods", "default", fields.Everything())

	reflector := cache.NewReflector(lwClient, &corev1.Pod{}, df, 0)

	go func() {
		reflector.Run(wait.NeverStop)
	}()

	for {
		_, _ = df.Pop(func(i interface{}) error {
			for _, delta := range i.(cache.Deltas) {
				pod := delta.Object.(*corev1.Pod)
				fmt.Println(delta.Type, pod.Name, pod.Namespace, pod.Status.Phase)

				switch delta.Type {
				case cache.Sync, cache.Added:
					store.Add(delta.Object)
				case cache.Updated:
					store.Update(delta.Object)
				case cache.Deleted:
					store.Delete(delta.Object)
				}
			}
			return nil
		})
	}
}
