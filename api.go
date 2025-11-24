// holds all the kubernetes api functions
package main

import (
	"bufio"
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	coreV1 "k8s.io/api/core/v1"
)

// lists all namespaces
func ListNamespaces(ctx context.Context, clientset *kubernetes.Clientset) ([]string, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if err != nil {
		return []string{}, err
	}

	names := make([]string, len(namespaces.Items))

	for i, ns := range namespaces.Items {
		names[i] = ns.Name
	}

	return names, nil
}

func ListPodsByNamespaces(ctx context.Context, clientset *kubernetes.Clientset, namespace string) ([]string, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return []string{}, err
	}

	podds := make([]string, len(pods.Items))

	for i, pod := range pods.Items {
		podds[i] = pod.Name
	}

	return podds, nil
}

func DeletePod(ctx context.Context, clientset *kubernetes.Clientset, namespace, pod string) error {
	return clientset.CoreV1().Pods(namespace).Delete(ctx, pod, v1.DeleteOptions{})
}

func ViewLogs(ctx context.Context, clientset *kubernetes.Clientset, namespace, pod string) (string, error) {
	count := int64(25)
	logOpts := coreV1.PodLogOptions{
		TailLines: &count,
	}
	req := clientset.CoreV1().Pods(namespace).GetLogs(pod, &logOpts)

	stream, err := req.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer stream.Close()

	scanner := bufio.NewScanner(stream)

	logs := ""
	for scanner.Scan() {
		logs += fmt.Sprintf("%s\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return logs, nil
}
