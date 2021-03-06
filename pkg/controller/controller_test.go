package controller

import (
	"context"
	"os"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/fairwindsops/astro/pkg/kube"
)

func TestCreateDeploymentController(t *testing.T) {
	kubeClient := kube.SetAndGetMock()
	DeploymentInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return kubeClient.Client.AppsV1().Deployments("").List(metav1.ListOptions{})
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return kubeClient.Client.AppsV1().Deployments("").Watch(metav1.ListOptions{})
			},
		},
		&appsv1.Deployment{},
		0,
		cache.Indexers{},
	)
	DeployWatcher := createController(kubeClient.Client, DeploymentInformer, "deployment")

	annotations := make(map[string]string, 1)
	annotations["test"] = "yup"
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "owned-namespace",
			Annotations: annotations,
		},
	}
	kubeClient.Client.CoreV1().Namespaces().Create(ns)
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
		},
	}
	kubeClient.Client.AppsV1().Deployments("owned-namespace").Create(deploy)

	assert.Implements(t, (*kubernetes.Interface)(nil), DeployWatcher.kubeClient, "")
	assert.Implements(t, (*cache.SharedIndexInformer)(nil), DeployWatcher.informer, "")
	assert.Implements(t, (*workqueue.RateLimitingInterface)(nil), DeployWatcher.wq, "")
}

func TestNewController(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	_, hook := test.NewNullLogger()
	log.AddHook(hook)
	kube.SetAndGetMock()
	os.Setenv("DEFINITIONS_PATH", "../config/test_conf.yml")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go New(ctx)

	time.Sleep(500 * time.Millisecond)
	var deployPass = false
	var namespacePass = false
	for _, log := range hook.AllEntries() {
		if deployPass && namespacePass {
			break
		}
		if log.Message == "Creating controller for resource type deployment" {
			deployPass = true
		}
		if log.Message == "Creating controller for resource type namespace" {
			namespacePass = true
		}
	}

	assert.Equal(t, true, deployPass, "Logging did not indicate that the deployment controller started.")
	assert.Equal(t, true, namespacePass, "Logging did not indicate that the namespace controller started.")
}
