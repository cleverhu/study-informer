package main

import (
	"fmt"
	
	"k8s.io/client-go/tools/cache"
)

type pod struct {
	Name  string
	Value float64
}

func newPod(name string, value float64) pod {
	return pod{Name: name, Value: value}
}

func main() {
	df := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{
		KeyFunction: func(obj interface{}) (string, error) {
			p := obj.(pod)
			return p.Name, nil
		},
	})
	p1 := newPod("pod1", 1)
	p2 := newPod("pod2", 2)
	p3 := newPod("pod3", 3)
	df.Add(p1)
	df.Add(p2)
	df.Add(p3)

	p1.Value = 1.1
	df.Update(p1)
	df.Delete(p1)

	fmt.Println(df.List())
	df.Pop(func(i interface{}) error {
		fmt.Printf("%T", i)
		deltas := i.(cache.Deltas)
		for _, delta := range deltas {
			fmt.Printf("%#v, Type: %v\n", delta.Object, delta.Type)
			switch delta.Type {
			case cache.Added:
				fmt.Println("执行新增操作")
			case cache.Updated:
				fmt.Println("执行更新操作")
			case cache.Deleted:
				fmt.Println("执行删除操作")
			}
		}
		return nil
	})

}
