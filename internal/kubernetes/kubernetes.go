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
	PodStates []PodState
}

type Cluster interface {
	IsExist() (bool, error)
	GetState() (ClusterInfo, error)

	Create() error
	Delete() error
	ScaleReplicas() error
}

type K8sCluster struct {
	Name         string
	Password     string
	ReplicaCount int
}

func (cls *K8sCluster) IsExist() (bool, error) {
	_, err := exec.Command("kubectl", "get", "sts", "-n", "mysql", cls.Name).Output()
	if err != nil {
		return false, err
	}
	return true, nil
}

func ParseClusterInfo(rawClusterInfo []byte) (error, ClusterInfo) {
	m := make(map[string]interface{})
	err := json.Unmarshal(rawClusterInfo, &m)
	if err != nil {
		return err, ClusterInfo{}
	}
	items, ok := m["items"].([]interface{})
	if !ok {
		return err, ClusterInfo{}
	}

	clusterInfo := ClusterInfo{}
	for i := 0; i < len(items); i++ {
		item, ok := items[i].(map[string]interface{})
		if !ok {
			return err, ClusterInfo{}
		}
		status, ok := item["status"].(map[string]interface{})
		if !ok {
			return err, ClusterInfo{}
		}
		phase, ok := status["phase"].(string)
		if !ok {
			return err, ClusterInfo{}
		}
		metadata, ok := item["metadata"].(map[string]interface{})
		if !ok {
			return err, ClusterInfo{}
		}
		name, ok := metadata["name"].(string)
		if !ok {
			return err, ClusterInfo{}
		}
		clusterInfo.PodStates = append(clusterInfo.PodStates, PodState{
			State:   phase,
			PodName: name,
		})
	}
	return err, clusterInfo
}

func (cls *K8sCluster) GetState() (ClusterInfo, error) {
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
	err, clusterInfo := ParseClusterInfo(output)
	return clusterInfo, err
}

func (cls *K8sCluster) Create() error {
	out, err := exec.Command(PATH+"create-cluster.sh",
		fmt.Sprintf("%s", cls.Name),
		fmt.Sprintf("%d", cls.ReplicaCount),
		fmt.Sprintf("%s", cls.Password),
	).Output()
	fmt.Println(string(out), err)
	return err
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
