package informer

import (
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"

	"github.com/cleverhu/study-informer/lib"
)

func Test_Example(t *testing.T) {
	cs := lib.InitClient()
	lw := cache.NewListWatchFromClient(cs.CoreV1().RESTClient(), "pods", "default", fields.Everything())
	i := NewInformer(lw, &corev1.Pod{}, &cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Println("Added: ", pod.Namespace, pod.Name, pod.Status.Phase)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			pod := newObj.(*corev1.Pod)
			fmt.Println("Updated: ", pod.Namespace, pod.Name, pod.Status.Phase)
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Println("Deleted: ", pod.Namespace, pod.Name, pod.Status.Phase)
		},
	})
	i.Run(wait.NeverStop)
}
