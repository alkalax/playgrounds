package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var server = "http://127.0.0.1:8001"

type K8sStatus struct {
	Message string `json:"message"`
}

func getNamespaces() []string {
	url := server + "/api/v1/namespaces"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("error: %s", body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var nsList NamespaceList
	if err := json.Unmarshal(body, &nsList); err != nil {
		panic(err)
	}

	namespaces := []string{}
	for _, ns := range nsList.Items {
		namespaces = append(namespaces, ns.Metadata.Name)
	}

	return namespaces
}

func getPods(namespace string) []string {
	url := fmt.Sprintf("%s/api/v1/namespaces/%s/pods", server, namespace)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("error: %s", body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var podList PodList
	if err := json.Unmarshal(body, &podList); err != nil {
		panic(err)
	}

	pods := []string{}
	for _, pod := range podList.Items {
		pods = append(pods, pod.Metadata.Name)
	}

	return pods
}

func getLogs(namespace, pod string) []string {
	url := fmt.Sprintf("%s/api/v1/namespaces/%s/pods/%s/log", server, namespace, pod)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var status K8sStatus
		if err := json.Unmarshal(body, &status); err != nil {
			panic(fmt.Sprintf("failed to parse error: %s", string(body)))
		}
		start := strings.Index(status.Message, "[")
		end := strings.Index(status.Message, "]")
		if start == -1 || end == -1 || start >= end {
			panic(fmt.Sprintf("error: %s", status.Message))
		}

		container := strings.Fields(status.Message[start+1 : end])[0]
		url := fmt.Sprintf("%s/api/v1/namespaces/%s/pods/%s/log?container=%s", server, namespace, pod, container)
		resp, err = http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return strings.Split(string(body), "\n")
}
