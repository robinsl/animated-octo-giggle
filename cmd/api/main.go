package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math/rand"
)

type Todo struct {
	gorm.Model
	Id        string `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(postgres.Open("postgresql://user:pass@localhost:5432/todo"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Todo{})

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, http://localhost:4321",
	}))

	todo := app.Group("/todo")

	todo.Get("/", func(c *fiber.Ctx) error {
		var todos []Todo
		db.Find(&todos)

		return c.JSON(todos)
	})

	todo.Get("/:id", func(c *fiber.Ctx) error {
		var todo Todo
		id := c.Params("id")
		db.First(&todo, id)

		return c.JSON(todo)
	})

	todo.Post("/", func(c *fiber.Ctx) error {
		todo := new(Todo)
		if err := c.BodyParser(todo); err != nil {
			return err
		}
		todo.Name = c.FormValue("name")
		db.Create(&todo)

		return c.JSON(todo)
	})

	todo.Put("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		var todo Todo
		if err := db.First(&todo, id).Error; err != nil {
			return c.Status(404).SendString("Todo item not found")
		}
		todo.Completed = !todo.Completed
		db.Model(&Todo{}).Where("ID = ?", id).Update("completed", todo.Completed)

		return c.JSON(todo)
	})

	todo.Delete("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		db.Delete(&Todo{}, id)
		return c.SendString("Todo item deleted: " + id)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
