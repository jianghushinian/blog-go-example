package main

import (
	"flag"
	"log/slog"
	"path/filepath"
	"time"

	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	nightwatch "github.com/jianghushinian/blog-go-example/nightwatch/internal"
	"github.com/jianghushinian/blog-go-example/nightwatch/pkg/db"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	var kubecfg *string
	if home := homedir.HomeDir(); home != "" {
		kubecfg = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "Optional absolute path to kubeconfig")
	} else {
		kubecfg = flag.String("kubeconfig", "", "Absolute path to kubeconfig")
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubecfg)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	config.QPS = 50
	config.Burst = 100
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	cfg := nightwatch.Config{
		MySQLOptions: &db.MySQLOptions{
			Host:                  "127.0.0.1:33306",
			Username:              "root",
			Password:              "nightwatch",
			Database:              "nightwatch",
			MaxIdleConnections:    100,
			MaxOpenConnections:    100,
			MaxConnectionLifeTime: time.Duration(10) * time.Second,
		},
		RedisOptions: &db.RedisOptions{
			Addr:         "127.0.0.1:36379",
			Username:     "",
			Password:     "nightwatch",
			Database:     0,
			MaxRetries:   3,
			MinIdleConns: 0,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolSize:     10,
		},
		Clientset: clientset,
	}

	nw, err := cfg.New()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	stopCh := genericapiserver.SetupSignalHandler()
	nw.Run(stopCh)
}
