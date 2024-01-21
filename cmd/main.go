package main

import (
	"dbaas-api/internal/kubernetes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func BasicAuthorizer(username, password string) bool {
	return true
}

func GetClusterFromBody(c *fiber.Ctx) (kubernetes.Cluster, error) {
	m := make(map[string]string)
	c.Accepts("json", "text")

	err := json.Unmarshal(c.Body(), &m)
	if err != nil {
		return nil, err
	}

	replicas, err := strconv.Atoi(m["replicas"])
	if err != nil {
		return nil, err
	}

	return &kubernetes.K8sCluster{
		Name:         m["name"],
		Password:     m["password"],
		ReplicaCount: replicas,
	}, nil
}

func main() {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Use(basicauth.New(basicauth.Config{
		Realm:      "Basic",
		Authorizer: BasicAuthorizer,
		Unauthorized: func(c *fiber.Ctx) error {
			return c.SendStatus(403)
		},
	}))

	app.Post("/api/state", func(c *fiber.Ctx) error {
		cls, err := GetClusterFromBody(c)
		if err != nil {
			return c.SendString(err.Error())
		}
		state, err := cls.GetState()
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.SendString(fmt.Sprint(state))
	})

	app.Post("/api/create", func(c *fiber.Ctx) error {
		cls, err := GetClusterFromBody(c)
		if err != nil {
			return c.SendString(err.Error())
		}
		go cls.Create()
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.SendString("Creating...")
	})

	app.Post("/api/delete", func(c *fiber.Ctx) error {
		cls, err := GetClusterFromBody(c)
		if err != nil {
			return c.SendString(err.Error())
		}
		go cls.Delete()
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.SendString("Deleting...")
	})

	app.Post("/api/update", func(c *fiber.Ctx) error {
		cls, err := GetClusterFromBody(c)
		if err != nil {
			return c.SendString(err.Error())
		}
		go cls.ScaleReplicas()
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.SendString("Updating...")
	})

	app.Listen(":8080")
}
