// Copyright 2018 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pingcap/tidb-operator/pkg/client/clientset/versioned"
	"github.com/pingcap/tidb-operator/pkg/features"
	"github.com/pingcap/tidb-operator/pkg/scheduler/server"
	"github.com/pingcap/tidb-operator/pkg/version"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"
)

var (
	printVersion bool
	port         int
)

func init() {
	klog.InitFlags(nil)
	flag.BoolVar(&printVersion, "V", false, "Show version and quit")
	flag.BoolVar(&printVersion, "version", false, "Show version and quit")
	flag.IntVar(&port, "port", 10262, "The port that the tidb scheduler's http service runs on (default 10262)")
	features.DefaultFeatureGate.AddFlag(flag.CommandLine)
	flag.Parse()
}

func main() {
	if printVersion {
		version.PrintVersionInfo()
		os.Exit(0)
	}
	version.LogVersionInfo()

	logs.InitLogs()
	defer logs.FlushLogs()

	flag.CommandLine.VisitAll(func(flag *flag.Flag) {
		klog.V(1).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
	})

	cfg, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatalf("failed to get config: %v", err)
	}
	kubeCli, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("failed to get kubernetes Clientset: %v", err)
	}
	cli, err := versioned.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("failed to create Clientset: %v", err)
	}

	go wait.Forever(func() {
		server.StartServer(kubeCli, cli, port)
	}, 5*time.Second)

	srv := http.Server{Addr: ":6060"}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		sig := <-sc
		klog.Infof("got signal %s to exit", sig)
		if err2 := srv.Shutdown(context.Background()); err2 != nil {
			klog.Fatal("fail to shutdown the HTTP server", err2)
		}
	}()

	if err = srv.ListenAndServe(); err != http.ErrServerClosed {
		klog.Fatal(err)
	}
	klog.Infof("tidb-scheduler exited")
}
