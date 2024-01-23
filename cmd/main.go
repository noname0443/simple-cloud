package main

import (
	"crypto/sha256"
	"dbaas-api/internal/kubernetes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/jmoiron/sqlx"
)

func SHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
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

type User struct {
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

func GetUserFromBody(c *fiber.Ctx) (User, error) {
	user := User{}
	c.Accepts("json", "text")

	err := json.Unmarshal(c.Body(), &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func main() {
	app := fiber.New(fiber.Config{
		Views: html.New("./static", ".html"),
	})

	connLine := fmt.Sprintf("%s:%s@(%s:%d)/", "admin", "admin", "192.168.1.65", 30255)
	db, err := sqlx.Open("mysql", connLine)
	if err != nil {
		panic(err)
	}

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})
	app.Get("/:page", func(c *fiber.Ctx) error {
		return c.Render(c.Params("page", "index"), fiber.Map{})
	})
	app.Post("/web/register", func(c *fiber.Ctx) error {
		email := c.FormValue("email")
		username := c.FormValue("username")
		password := c.FormValue("password")
		_, err := db.Query("INSERT INTO metadb.users (email, username, password) VALUES (?, ?, ?);", email, username, SHA256([]byte(password)))
		log.Print(err)
		if err != nil {
			return c.SendStatus(500)
		}
		c.Cookie(&fiber.Cookie{
			Name:  "AUTH",
			Value: username + ":" + SHA256([]byte(password)),
		})
		return c.SendStatus(200)
	})

	validateToken := func(c *fiber.Ctx, key string) (bool, error) {
		data := strings.Split(key, ":")
		username, password := data[0], data[1]
		var user User
		err := db.Select(&user, "SELECT username, password FROM metadb.users WHERE username = ? AND password = ?", username, password)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	api := app.Group("/api", keyauth.New(keyauth.Config{
		KeyLookup: "cookie:AUTH",
		Validator: validateToken,
	}))

	api.Post("/api/state", func(c *fiber.Ctx) error {
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

	api.Post("/api/create", func(c *fiber.Ctx) error {
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

	api.Post("/api/delete", func(c *fiber.Ctx) error {
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

	api.Post("/api/update", func(c *fiber.Ctx) error {
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
