package main

import (
	"crypto/sha256"
	"dbaas-api/internal/kubernetes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/handlebars/v2"
	"github.com/jmoiron/sqlx"
)

func SHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

func GetClusterFromBody(c *fiber.Ctx) (*kubernetes.K8sCluster, error) {
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
		Username:     m["username"],
		Name:         m["name"],
		Password:     m["password"],
		DatabaseName: m["database"],
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
	IP := flag.String("ip", "0.0.0.0", "metadb database IP")
	PORT := flag.Int("port", 0, "metadb database port")
	flag.Parse()
	log.Println("metadb", *IP, *PORT)

	app := fiber.New(fiber.Config{
		Views: handlebars.New("./static", ".hbs"),
	})

	app.Static("/styles", "./static/styles")
	app.Static("/js", "./static/js")

	connLine := fmt.Sprintf("%s:%s@(%s:%d)/", "admin", "admin", *IP, *PORT)
	db, err := sqlx.Connect("mysql", connLine)
	if err != nil {
		panic(err)
	}

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})
	app.Get("/states", func(c *fiber.Ctx) error {
		states, err := GetClusterInfoStates(db, c)
		if err != nil {
			log.Println(err)
			return c.SendStatus(500)
		}
		return c.Render(c.Params("page", "states"), fiber.Map{
			"states":   states,
			"password": GetPassword(c),
			"username": GetUsername(c),
			"IP":       IP,
		})
	})
	app.Get("/:page", func(c *fiber.Ctx) error {
		return c.Render(c.Params("page", "index"), fiber.Map{})
	})
	app.Post("/web/register", func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		rows, err := db.Query("SELECT * FROM metadb.users WHERE username = ? AND password = ?;", username, SHA256([]byte(password)))
		if !rows.Next() {
			_, err = db.Query("INSERT INTO metadb.users (username, password) VALUES (?, ?);", username, SHA256([]byte(password)))
			log.Println(err)
			if err != nil {
				return c.SendStatus(500)
			}
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
		err := db.Get(&user, "SELECT username, password FROM metadb.users WHERE username = ? AND password = ?", username, password)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	api := app.Group("/api", keyauth.New(keyauth.Config{
		KeyLookup: "cookie:AUTH",
		Validator: validateToken,
	}))

	api.Get("/states", func(c *fiber.Ctx) error {
		states, err := GetClusterInfoStates(db, c)
		if err != nil {
			return c.SendString(err.Error())
		}
		jsonData, err := json.Marshal(states)
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.Send(jsonData)
	})

	api.Post("/create", func(c *fiber.Ctx) error {
		cls := kubernetes.K8sCluster{}
		err := json.Unmarshal(c.Body(), &cls)
		if err != nil {
			return c.SendString(err.Error())
		}
		username := GetUsername(c)
		if err != nil {
			return c.SendString(err.Error())
		}
		go func() {
			_, err = db.Query("INSERT INTO metadb.clusters (username, cluster_name, cluster_username, cluster_password) VALUES (?, ?, ?, ?);", username, cls.Name, cls.Username, cls.Password)
			if err != nil {
				log.Println(err)
			}
			err := cls.Create()
			if err != nil {
				log.Println(err)
				_, err = db.Exec("DELETE FROM metadb.clusters WHERE username = ? AND cluster_name = ?;", username, cls.Name)
				if err != nil {
					log.Println(err)
				}
			}
		}()
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.SendString("Creating...")
	})

	api.Post("/delete", func(c *fiber.Ctx) error {
		cls := kubernetes.K8sCluster{}
		err := json.Unmarshal(c.Body(), &cls)
		if err != nil {
			return c.SendString(err.Error())
		}
		username := GetUsername(c)
		go func() {
			_, err = db.Exec("DELETE FROM metadb.clusters WHERE username = ? AND cluster_name = ?;", username, cls.Name)
			if err != nil {
				log.Println(err)
			}
			cls.Delete()
		}()
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.SendString("Deleting...")
	})

	api.Post("/update", func(c *fiber.Ctx) error {
		cls := kubernetes.K8sCluster{}
		err := json.Unmarshal(c.Body(), &cls)
		if err != nil {
			return c.SendString(err.Error())
		}
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
