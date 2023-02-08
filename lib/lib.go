package lib

import (
	"log"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	oc = sync.Once{}
	cs kubernetes.Interface
)

func InitClient() kubernetes.Interface {
	if cs != nil {
		return cs
	}
	oc.Do(func() {
		config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			log.Fatalln(err)
		}

		cs, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalln(err)
		}
	})
	return cs
}
