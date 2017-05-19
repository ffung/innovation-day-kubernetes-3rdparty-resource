package main

var (
	apiHost                   = "http://127.0.0.1:8001"
	environmentsEndpoint      = "/apis/stable.xebia.com/v1/namespaces/default/environment"
	environmentsWatchEndpoint = "/apis/stable.xebia.com/v1/namespaces/default/environment?watch=true"
)

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

func getEnvironments() ([]Environments, error) {
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
