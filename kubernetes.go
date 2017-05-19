package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ericchiang/k8s"
	"github.com/ericchiang/k8s/api/v1"
	metav1 "github.com/ericchiang/k8s/apis/meta/v1"
	"log"
	"net/http"
	"time"
)

var (
	apiHost                   = "http://127.0.0.1:8080"
	environmentsEndpoint      = "/apis/stable.xebia.com/v1/namespaces/default/environments"
	environmentsWatchEndpoint = "/apis/stable.xebia.com/v1/namespaces/default/environments?watch=true"
)

type EnvironmentEvent struct {
	Type   string      `json:"type"`
	Object Environment `json:"object"`
}

type Environment struct {
	ApiVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   map[string]string `json:"metadata"`
	Spec       EnvironmentSpec   `json:"spec"`
}

type EnvironmentSpec struct {
	EnvironmentNamespace string `json:"environment-namespace"`
}

type EnvironmentList struct {
	ApiVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   map[string]string `json:"metadata"`
	Items      []Environment     `json:"items"`
}

func getEnvironments() ([]Environment, error) {
	var resp *http.Response
	var err error
	for {
		resp, err = http.Get(apiHost + environmentsEndpoint)
		if err != nil {
			log.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	var envList EnvironmentList
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&envList)
	if err != nil {
		return nil, err
	}

	return envList.Items, nil
}

func monitorEnvironmentEvents() (<-chan EnvironmentEvent, <-chan error) {
	events := make(chan EnvironmentEvent)
	errc := make(chan error, 1)
	go func() {
		for {
			resp, err := http.Get(apiHost + environmentsWatchEndpoint)
			if err != nil {
				errc <- err
				time.Sleep(5 * time.Second)
				continue
			}
			if resp.StatusCode != 200 {
				errc <- errors.New("Invalid status code: " + resp.Status)
				time.Sleep(5 * time.Second)
				continue
			}

			decoder := json.NewDecoder(resp.Body)
			for {
				var event EnvironmentEvent
				err = decoder.Decode(&event)
				if err != nil {
					errc <- err
					break
				}
				events <- event
			}
		}
	}()

	return events, errc
}

func createNamespace(client *k8s.Client, name string) error {
	ns := &v1.Namespace{
		Metadata: &metav1.ObjectMeta{
			Name: &name,
		},
	}
	_, err := client.CoreV1().CreateNamespace(context.TODO(), ns)

	if err != nil {
		log.Printf("Created namespace: %s", name)
	} else {
		log.Printf("Creating namespace failed: %s, reason: %v", name, err)
	}

	return err
}

func deleteNamespace(client *k8s.Client, name string) error {
	err := client.CoreV1().DeleteNamespace(context.TODO(), name)
	if err != nil {
		log.Printf("Deleted namespace: %s", name)
	} else {
		log.Printf("Deleting namespace failed: %s, reason: %v", name, err)
	}
	return err
}
