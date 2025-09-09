package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task-manager-app/config"
	"task-manager-app/controllers"
	"task-manager-app/models"
	"task-manager-app/utils"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/v2/bson"
)
var _ = Describe("AuthController", func (){
	var app *fiber.App
	var loginRequest  struct{
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	BeforeEach(func (){
		app = fiber.New()
		config.ConnectDB()
		config.DB.Collection("users").DeleteMany(context.Background(), bson.M{})
		hashedPassword, _ := utils.HashPassword("password123")
		users := models.User{
			Username:  "testuser",
			Email: "test@example.com",
			Password: hashedPassword,
		}
		_,err := config.DB.Collection("users").InsertOne(context.Background(),users)
		Expect(err).NotTo(HaveOccurred())

		app.Post("auth/login", controllers.Login)
	})

	Context("Saat User melakukan Login", func(){
		It("Harus mengembalikan Token JWT jika kredensial Benar", func (){
			loginRequest.Email = "test@example.com"
			loginRequest.Password = "password123"

			body, _ := json.Marshal(loginRequest)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp,err := app.Test(req, -1)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			
			var responBody struct{
				Token string `json:"token"`
			}
			json.NewDecoder(req.Body).Decode(&responBody)
			Expect(responBody.Token).ToNot(BeEmpty())
		})
		It("Harus mengembalikan error jika password salah", func() {
			loginRequest.Email = "test@example.com"
			loginRequest.Password = "salahpassword"

			body, _ := json.Marshal(loginRequest)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
		})
	})
})