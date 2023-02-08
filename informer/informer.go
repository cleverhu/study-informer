package informer

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

type Informer struct {
	lw           *cache.ListWatch
	expectedType interface{}
	handler      cache.ResourceEventHandler

	store     cache.Store
	fifo      *cache.DeltaFIFO
	reflector *cache.Reflector
}

func NewInformer(lw *cache.ListWatch, expectedType interface{}, handler cache.ResourceEventHandler) *Informer {
	store := cache.NewStore(cache.MetaNamespaceKeyFunc)
	fifo := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{
		KeyFunction:  cache.MetaNamespaceKeyFunc,
		KnownObjects: store,
	})
	reflector := cache.NewReflector(lw, &corev1.Pod{}, fifo, 0)
	return &Informer{lw: lw, expectedType: expectedType, handler: handler, store: store, fifo: fifo, reflector: reflector}
}

func (i *Informer) Run(stopCh <-chan struct{}) {
	go func() {
		i.reflector.Run(stopCh)
	}()

	for {
		_, _ = i.fifo.Pop(func(data interface{}) error {
			for _, delta := range data.(cache.Deltas) {
				switch delta.Type {
				case cache.Sync, cache.Added:
					i.store.Add(delta.Object)
					i.handler.OnAdd(delta.Object)
				case cache.Updated:
					if item, exist, err := i.store.Get(delta.Object); err == nil && exist {
						i.store.Update(delta.Object)
						i.handler.OnUpdate(item, delta.Object)
					}
				case cache.Deleted:
					i.store.Delete(delta.Object)
					i.handler.OnDelete(delta.Object)
				}
			}
			return nil
		})
	}
}
