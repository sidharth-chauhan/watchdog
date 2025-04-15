package tests

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	namespace = "watchdog-ns"
	timeout   = 60 * time.Second
	interval  = 5 * time.Second
)

func TestKubernetesManifests(t *testing.T) {
	// Get kubeconfig
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatalf("Error building kubeconfig: %v", err)
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Error creating clientset: %v", err)
	}

	// Test 1: Validate deployment configuration
	t.Run("Validate deployment configuration", func(t *testing.T) {
		deploymentPath := filepath.Join("..", "deployment.yaml")
		if _, err := os.Stat(deploymentPath); os.IsNotExist(err) {
			t.Fatalf("Deployment file not found: %v", err)
		}
	})

	// Test 2: Validate service configuration
	t.Run("Validate service configuration", func(t *testing.T) {
		servicePath := filepath.Join("..", "service.yaml")
		if _, err := os.Stat(servicePath); os.IsNotExist(err) {
			t.Fatalf("Service file not found: %v", err)
		}
	})

	// Test 3: Validate configmap configuration
	t.Run("Validate configmap configuration", func(t *testing.T) {
		configmapPath := filepath.Join("..", "configmap.yaml")
		if _, err := os.Stat(configmapPath); os.IsNotExist(err) {
			t.Fatalf("ConfigMap file not found: %v", err)
		}
	})

	// Test 4: Check pod status
	t.Run("Check pod status", func(t *testing.T) {
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: "app=watchdog",
		})
		if err != nil {
			t.Fatalf("Error listing pods: %v", err)
		}

		if len(pods.Items) == 0 {
			t.Error("No pods found in namespace")
			return
		}

		pod := pods.Items[0]
		if pod.Status.Phase != "Running" {
			t.Errorf("Expected pod status Running, got %s", pod.Status.Phase)
		}

		if pod.Status.ContainerStatuses[0].Ready != true {
			t.Error("Container is not ready")
		}
	})

	// Test 5: Check service
	t.Run("Check service", func(t *testing.T) {
		service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), "watchdog-service", metav1.GetOptions{})
		if err != nil {
			t.Fatalf("Error getting service: %v", err)
		}

		if service.Spec.Type != "NodePort" {
			t.Errorf("Expected service type NodePort, got %s", service.Spec.Type)
		}

		if len(service.Spec.Ports) == 0 {
			t.Error("No ports configured in service")
		} else if service.Spec.Ports[0].Port != 8080 {
			t.Errorf("Expected port 8080, got %d", service.Spec.Ports[0].Port)
		}
	})

	// Test 6: Check deployment
	t.Run("Check deployment", func(t *testing.T) {
		deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), "watchdog", metav1.GetOptions{})
		if err != nil {
			t.Fatalf("Error getting deployment: %v", err)
		}

		if *deployment.Spec.Replicas != 1 {
			t.Errorf("Expected 1 replica, got %d", *deployment.Spec.Replicas)
		}

		if deployment.Status.AvailableReplicas != 1 {
			t.Errorf("Expected 1 available replica, got %d", deployment.Status.AvailableReplicas)
		}
	})

	// Test 7: Check ConfigMap
	t.Run("Check ConfigMap", func(t *testing.T) {
		configmap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), "watchdog-config", metav1.GetOptions{})
		if err != nil {
			t.Fatalf("Error getting ConfigMap: %v", err)
		}

		if _, exists := configmap.Data["oba_server_config.json"]; !exists {
			t.Error("ConfigMap does not contain oba_server_config.json")
		}
	})
}
