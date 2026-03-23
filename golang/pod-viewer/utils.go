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
