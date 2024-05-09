package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DaffaJatmiko/go-fiber-postgres/models"
	"github.com/DaffaJatmiko/go-fiber-postgres/storage"
	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author string `json:"author"`
	Title string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) {
	book := &Book{}

	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"massage": err.Error()})
			return 
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"massage": err.Error()})
			return 
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Success create a book"})

}

func (r *Repository) DeleteBook(context *fiber.Ctx) {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return 
	}

	err := r.DB.Delete(bookModel, id).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Could not delete book",
		})
		return 
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message" : "success delete book",
	})

}

func (r *Repository) GetBooks(context *fiber.Ctx) {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": "Book not found"})
			return 
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "book fetched successfully",
			"book": bookModels,
		})

}

func (r *Repository) GetBookByID(context *fiber.Ctx) {
	bookModel := &models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return 
	}

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the book"})
			return
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "book id fetched successfully",
			"book": bookModel,
		})

}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User: os.Getenv("DB_USER"),
		DBName: os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("failed to load database: ", err)
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migrage db")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}