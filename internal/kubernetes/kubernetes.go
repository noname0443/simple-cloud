package kubernetes

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

const PATH = "./"

type PodState struct {
	PodName string
	State   string
}

type ClusterInfo struct {
	PodStates   []PodState
	ClusterName string
	Port        int
}

type K8sCluster struct {
	Name         string `json:"name,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	DatabaseName string `json:"database,omitempty"`
	ReplicaCount int    `json:"replicas,omitempty"`
}

func ParsePort(rawPortInfo []byte) (int, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(rawPortInfo, &m)
	if err != nil {
		return -1, err
	}
	spec, ok := m["spec"].(map[string]interface{})
	if !ok {
		return -2, err
	}

	ports, ok := spec["ports"].([]interface{})
	if !ok || len(ports) == 0 {
		return -3, err
	}

	firstPort, ok := ports[0].(map[string]interface{})
	if !ok {
		return -4, err
	}

	floatNodePort, ok := firstPort["nodePort"].(float64)
	nodePort := int(floatNodePort)
	if !ok {
		return -5, err
	}
	return nodePort, err
}

func ParsePodStates(rawClusterInfo []byte) ([]PodState, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(rawClusterInfo, &m)
	if err != nil {
		return []PodState{}, err
	}
	items, ok := m["items"].([]interface{})
	if !ok {
		return []PodState{}, err
	}

	podStates := make([]PodState, len(items))
	for i := 0; i < len(items); i++ {
		item, ok := items[i].(map[string]interface{})
		if !ok {
			return []PodState{}, err
		}
		status, ok := item["status"].(map[string]interface{})
		if !ok {
			return []PodState{}, err
		}
		phase, ok := status["phase"].(string)
		if !ok {
			return []PodState{}, err
		}
		metadata, ok := item["metadata"].(map[string]interface{})
		if !ok {
			return []PodState{}, err
		}
		name, ok := metadata["name"].(string)
		if !ok {
			return []PodState{}, err
		}
		podStates[i] = PodState{
			State:   phase,
			PodName: name,
		}
	}
	return podStates, err
}

func (cls *K8sCluster) GetState() (ClusterInfo, error) {
	_, err := exec.Command(PATH+"forward-svc.sh",
		fmt.Sprintf("%s", cls.Name),
	).Output()
	if err != nil {
		return ClusterInfo{}, err
	}

	output, err := exec.Command(
		"kubectl",
		"get",
		"pods",
		"-n",
		"mysql",
		"-l",
		fmt.Sprintf("app.kubernetes.io/instance=%s", cls.Name),
		"-o",
		"json",
	).Output()
	if err != nil {
		return ClusterInfo{}, err
	}
	podStates, err := ParsePodStates(output)

	output, err = exec.Command(
		"kubectl",
		"get",
		"svc",
		"-n",
		"services",
		cls.Name,
		"-o",
		"json",
	).Output()
	if err != nil {
		return ClusterInfo{}, err
	}
	port, err := ParsePort(output)

	return ClusterInfo{
		PodStates:   podStates,
		Port:        port,
		ClusterName: cls.Name,
	}, err
}

func (cls *K8sCluster) Create() error {
	out, err := exec.Command(PATH+"create-cluster.sh",
		fmt.Sprintf("%s", cls.Name),
		fmt.Sprintf("%d", cls.ReplicaCount),
		fmt.Sprintf("%s", cls.Password),
		fmt.Sprintf("%s", cls.Username),
		fmt.Sprintf("%s", cls.DatabaseName),
	).Output()
	if err != nil {
		return fmt.Errorf("%s:%s", out, err.Error())
	}
	return nil
}

func (cls *K8sCluster) Delete() error {
	_, err := exec.Command(PATH+"delete-cluster.sh",
		fmt.Sprintf("%s", cls.Name),
	).Output()
	return err
}

func (cls *K8sCluster) ScaleReplicas() error {
	_, err := exec.Command(PATH+"update-cluster.sh",
		fmt.Sprintf("%s", cls.Name),
		fmt.Sprintf("%d", cls.ReplicaCount),
	).Output()
	return err
}
