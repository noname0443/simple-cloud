package main

import (
	"dbaas-api/internal/kubernetes"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func GetClusterInfoStates(db *sqlx.DB, c *fiber.Ctx) ([]kubernetes.ClusterInfo, error) {
	clusters := []string{}

	data := strings.Split(c.Cookies("AUTH"), ":")
	username := data[0]

	err := db.Select(&clusters, "SELECT cluster_name FROM metadb.clusters WHERE username = ?", username)
	if err != nil {
		return nil, err
	}

	var states []kubernetes.ClusterInfo
	for i := 0; i < len(clusters); i++ {
		cls := kubernetes.K8sCluster{
			Name: clusters[i],
		}
		state, err := cls.GetState()
		if err != nil {
			return nil, err
		}
		states = append(states, state)
	}
	return states, err
}

func GetPassword(c *fiber.Ctx) string {
	data := strings.Split(c.Cookies("AUTH"), ":")
	if len(data) < 2 {
		return ""
	}
	return data[1]
}

func GetUsername(c *fiber.Ctx) string {
	data := strings.Split(c.Cookies("AUTH"), ":")
	if len(data) < 2 {
		return ""
	}
	return data[0]
}
