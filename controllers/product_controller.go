package controllers

import (
	"context"
	"fmt"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"gin-mongo-api/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var productCollection *mongo.Collection = configs.GetCollection(configs.DB, "products")

func CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var product models.Product
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&product); err != nil {
			fmt.Printf("1\n")
			c.JSON(http.StatusBadRequest, responses.ProductResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&product); validationErr != nil {
			fmt.Printf("2\n")
			c.JSON(http.StatusBadRequest, responses.ProductResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newProduct := models.Product{
			Id:       primitive.NewObjectID(),
			Name:     product.Name,
			SKU:      product.SKU,
			Quantity: product.Quantity,
		}

		result, err := productCollection.InsertOne(ctx, newProduct)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProductResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.ProductResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetAProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		productId := c.Param("productId")
		var product models.Product
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(productId)

		err := productCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProductResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.ProductResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": product}})
	}
}

func EditAProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		productId := c.Param("productId")
		var product models.Product
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(productId)

		//validate the request body
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, responses.ProductResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&product); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.ProductResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"name": product.Name, "SKU": product.SKU, "quantity": product.Quantity}
		result, err := productCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProductResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated product details
		var updatedProduct models.Product
		if result.MatchedCount == 1 {
			err := productCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedProduct)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.ProductResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.ProductResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedProduct}})
	}
}

func DeleteAProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		productId := c.Param("productId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(productId)

		result, err := productCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProductResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.ProductResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Product with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.ProductResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Product successfully deleted!"}},
		)
	}
}

func GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var products []models.Product
		defer cancel()

		results, err := productCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProductResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleProduct models.Product
			if err = results.Decode(&singleProduct); err != nil {
				c.JSON(http.StatusInternalServerError, responses.ProductResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			products = append(products, singleProduct)
		}

		c.JSON(http.StatusOK,
			responses.ProductResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": products}},
		)
	}
}
