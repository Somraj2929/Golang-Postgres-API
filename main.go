package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/somraj/go-fiber/postgress/models"
	"github.com/somraj/go-fiber/postgress/storage"
	"gorm.io/gorm"
)

type Book struct{
	Author 		string		`json:"author"`
	Title		string		`json:"title"`
	Publisher	string		`json:"publisher"`
}

//defining datatype for database
type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request Failed"})
		return err
	}

	err = r.DB.Create(&book).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Not able to create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "Book created successfully"})
	return nil
}

func (r *Repository) GetBook(context *fiber.Ctx) error{
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"Unable to find Books"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message":"Book fetched successfully",
					"data" : bookModels,
				})
		return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := &[]models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(bookModel, id).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err
	}
	fmt.Println("Deletion API Called")
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book delete successfully",
	})
	return nil
}


func (r *Repository) GetBookByID(context *fiber.Ctx) error{
	id := context.Params("id")
	bookModel := &models.Books{}
	if id == ""{
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message":"Id cannot be empty",
		})
		return nil
	}

	fmt.Println("The ID is", id)

	err := r.DB.Where("id = ?", id).First(bookModel).Error 
	if err != nil{
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message":"Failed to get book.",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message":"Book fetched successfully",
		"data" : bookModel,
	})
	return nil
}



//calling own method 
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBook)
}


func main() {
	err := godotenv.Load(os.ExpandEnv("C:/Users/ASUS/Desktop/Golang/Golang-Postgress-API/.env"))
	
	if err != nil {
		log.Fatal(err)
	}


	config := &storage.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User: os.Getenv("DB_USER"),
		SSLMode: os.Getenv("DB_SSLMODE"),
		DBName: os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)
	if err != nil{
		log.Fatal("Could not load database")
	}

	err = models.MigrateBooks(db)
	if err != nil{
		log.Fatal("Could not migrate db")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}