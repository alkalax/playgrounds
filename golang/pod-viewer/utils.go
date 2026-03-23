package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var server = "http://127.0.0.1:8001"

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
