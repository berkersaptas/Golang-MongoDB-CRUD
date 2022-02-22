package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	_ "github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var humansCollection = db().Database("Animals").Collection("Human")

type Human struct {
	//ID      primitive.ObjectID `bson:"_id,omitempty"`
	No      string `bson:"no,omitempty", valid:"required"`
	Name    string `bson:"name,omitempty" valid:"required"`
	Surname string `bson:"surname,omitempty" valid:"required"`
	Gender  string `bson:"gender,omitempty" valid:"required"`
	Age     string `bson:"age,omitempty" valid:"required"`
	Active  *bool  `bson:"active,omitempty"`
}

type ResponseStruct struct {
	Error   bool
	Message string
}

//db connection
func db() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func adduser(c echo.Context) error {
	human := &Human{}
	responceMessage := &ResponseStruct{}

	//read body raw
	body, err := ioutil.ReadAll(c.Request().Body)

	if err != nil {
		responceMessage.Error = true
		responceMessage.Message = "Req Body is not read"
		return c.JSON(http.StatusBadRequest, responceMessage)
	}
	//jsondecode
	err = json.Unmarshal(body, &human)
	if err != nil {
		responceMessage.Error = true
		responceMessage.Message = "Json decode error"
		return c.JSON(http.StatusBadRequest, responceMessage)
	}
	//check request model
	_, err = valid.ValidateStruct(human)
	if err != nil {
		responceMessage.Error = true
		responceMessage.Message = "Req json not match struct validate"
		return c.JSON(http.StatusBadRequest, responceMessage)
	}
	//db prossess
	_, err = humansCollection.InsertOne(context.TODO(), human)
	if err != nil {
		responceMessage.Error = true
		responceMessage.Message = fmt.Sprintf("DB insert not succesfull err is: %d", err)
		return c.JSON(http.StatusBadRequest, responceMessage)
	}
	responceMessage.Error = false
	responceMessage.Message = "DB add user succesfull"
	return c.JSON(http.StatusOK, responceMessage)
}

func readspesificuser(c echo.Context) error {
	//read get request param
	no := c.QueryParam("no")
	var result primitive.M
	//db prossess search user no
	err := humansCollection.FindOne(context.TODO(), bson.D{{"no", no}}).Decode(&result)
	//return http 500
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}
func readusers(c echo.Context) error {

	//db prossess
	cursor, err := humansCollection.Find(
		context.TODO(), bson.D{{}})
	if err != nil {
		return err
	}
	var humans []Human

	if err = cursor.All(context.TODO(), &humans); err != nil {
		return err
	}

	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, humans)
}

func updateuser(c echo.Context) error {
	//read get request param
	no := c.QueryParam("no")

	human := &Human{}
	responceMessage := &ResponseStruct{}

	//read body raw
	body, err := ioutil.ReadAll(c.Request().Body)

	if err != nil {
		return err
	}

	//jsondecode
	err = json.Unmarshal(body, &human)
	if err != nil {
		return err
	}

	//db prossess
	filter := bson.D{{"no", no}}
	update := bson.D{{"$set", human}}
	updateResult, err := humansCollection.UpdateOne(context.TODO(), filter, update)
	if updateResult.ModifiedCount != 1 {
		responceMessage.Error = true
		responceMessage.Message = "Modified error"
		return c.JSON(http.StatusBadRequest, responceMessage)
	}
	if err != nil {
		responceMessage.Error = true
		responceMessage.Message = "DB writing error"
		return c.JSON(http.StatusBadRequest, responceMessage)
	}
	responceMessage.Error = false
	responceMessage.Message = "Modified Successfully"
	return c.JSON(http.StatusOK, responceMessage)
}

func deleteuser(c echo.Context) error {
	no := c.QueryParam("no")
	user := bson.D{{"no", no}}
	responceMessage := &ResponseStruct{}

	// db prossess
	deleteResult, err := humansCollection.DeleteOne(context.TODO(), user)
	if err != nil {
		responceMessage.Error = true
		responceMessage.Message = "DB error"
		return c.JSON(http.StatusBadRequest, responceMessage)
	}
	if deleteResult.DeletedCount != 1 {
		if err != nil {
			responceMessage.Error = true
			responceMessage.Message = "User Not Deleted"
			return c.JSON(http.StatusBadRequest, responceMessage)
		}
	}
	responceMessage.Error = false
	responceMessage.Message = "User Deleted"
	return c.JSON(http.StatusOK, responceMessage)
}

func main() {
	e := echo.New()
	e.POST("/adduser", adduser)
	e.GET("/readspesificuser", readspesificuser)
	e.GET("/readusers", readusers)
	e.POST("/updateuser", updateuser)
	e.GET("/deleteuser", deleteuser)
	e.Logger.Fatal(e.Start(":1323"))
}
