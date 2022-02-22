### Golang Mongodb CRUD operation

In this project, I aimed to experiment and learn.

User information transmitted to the server as json in the project is transmitted to mongodb for registration.

## Mongodb install
[First install, let's follow the official document for mongodb installation](https://docs.mongodb.com/manual/installation/)

## Golang install
[Secound install is go lang](https://go.dev/doc/install)

## Postman install
We need to install an application that can make http requests. For this, it can be used in CURL via terminal, but I preferred Postman.[postman official install document](https://learning.postman.com/docs/getting-started/installation-and-updates/)

# For the DB connection, we should create a function as follows.

We can make the mongo db connection with the standard community version installed as follows.

```sh
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

```


# I created the following type of collection and document on mongodb.
 ```sh
 humansCollection = db().Database("Animals").Collection("Human")
 ```

# CRUD operation

Create , Read , Update and Delete functions in the 