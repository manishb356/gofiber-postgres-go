package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/manishb356/gofiber-postgres-go/models"
	"github.com/manishb356/gofiber-postgres-go/storage"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetBooks(ctx *fiber.Ctx) error {
	bookModels := &[]models.Book{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not get books"})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "Books fetched successfully", "data": bookModels})

	return nil
}

func (r *Repository) GetBook(ctx *fiber.Ctx) error {
	bookModel := &models.Book{}
	id := ctx.Params("id")

	if id == "" {
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "ID is required"})
		return nil
	}

	err := r.DB.Find(bookModel, id).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not get book with id"})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book fetched successfully", "data": bookModel})

	return nil
}

func (r *Repository) CreateBook(ctx *fiber.Ctx) error {
	book := Book{}

	err := ctx.BodyParser(&book)
	if err != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Request failed"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not create book"})
		return err
	}

	ctx.Status(http.StatusCreated).JSON(&fiber.Map{"message": "Book added"})

	return nil
}

func (r *Repository) DeleteBook(ctx *fiber.Ctx) error {
	bookModel := &models.Book{}
	id := ctx.Params("id")

	if id == "" {
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "ID is required"})
		return nil
	}

	err := r.DB.Delete(bookModel, id).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not delete book"})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book deleted successfully", "data": bookModel})

	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/books", r.GetBooks)
	api.Get("/get_book/:id", r.GetBook)
	api.Post("/create_book", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Something went wrong: ", err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not load database, ", err)
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Could not migrate db, ", err)
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
