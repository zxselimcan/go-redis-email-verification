package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var (
	redisClient *redis.Client
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	redisClient = ConnectRedis()
}

func main() {
	initRoutes()
}

func initRoutes() {
	app := fiber.New()
	app.Get("/send-verification-email", startVerification) // Don't use it on prod. You should start verification inside your signup route. Also don't use GET method for this.
	app.Get("/verify-email", verifyEmail)                  // send a post request to this. use insominia or postman to do it. or do it on clientside. If you're on production, don't send this link in mail or don't use GET method for verification.
	app.Listen(":3000")

}

func ConnectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		DB:   0,
	})
}

func startVerification(c *fiber.Ctx) error {
	email := c.Query("email")
	verificationLink, err := SendVerificationMail(email)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{
			"message": "lmao its not working",
		})
	}

	fmt.Printf("Here is your link: %v", verificationLink)

	return c.JSON(fiber.Map{
		"message": "Sent!",
	})
}

func verifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	email, err := redisClient.Get(context.Background(), token).Result()
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Error"})
	}
	fmt.Println(email)
	fmt.Println("Now you can update emailVerified column to true")

	_, err = redisClient.Del(context.Background(), token).Result()
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Error"})
	}

	return c.JSON(fiber.Map{
		"message": "verified",
	})
}
