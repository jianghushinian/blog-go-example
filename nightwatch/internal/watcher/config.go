package watcher

import (
	"k8s.io/client-go/kubernetes"

	"github.com/jianghushinian/blog-go-example/nightwatch/pkg/store"
)

type Config struct {
	Store store.IStore

	Clientset kubernetes.Interface
}
