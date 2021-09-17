# qa-api
Q&amp;A Basic API developed with GO and GOKIT

The technologies implemented in this API are:
- Go
- Gorilla/mux
- gRPC
- Gokit
- Docker
- Testify
- MondoDB

# System requirements

This are the minimum system requirements that will be needed to execute the API
- Go
- Docker
- MongoDB
- MockGen (recommended for testing purposes)
- Protoc (used for protobuff files generation)

# Get going

Once you have the project in your machine and the minimun system requirments, you will need to complete the following steps:

- On the base project directory, execute `go mod tidy` and then `go mod vendor` to donwload all the dependencies of the project

- By default the API start with its own Mongo database manage by the docker compose file configuration, if you're using a remote mongoDB host, make sure that change the database URI variable (MONGODB_URI) on the .env file
- Run the command `docker compose build` to build the docker image

- Once the image if ready, run the command `docker compose up` to start the server

- OPTIONAL: if you want the server running in detached moded, run the command `docker compose up -d`

NOTE: You can test the REST endpoints with the request.http file, if you wanna test the gRPC endpoints, open a new terminal (while the app is running) and run the command `go run cmd/questionary/client/grpc/main.go`

And thats it, you are ready to GO :)



