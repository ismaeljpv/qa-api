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

# Get going

Once you have the project in your machine and the minimun system requirments, you will need to complete the folloring steps

- Execute `go mod tidy` and then `go mod vendor` to donwload all the dependencies of the project
- Make sure that add the right database URI on the .env file
- On the base project directory, run the command `docker compose build` to build the docker image
- Once the image if ready, run the command `docker compose up` to start the server. 
- OPTIONAL: if you want the server running in detached moded, run the command `docker compose up -d`.

If your running the mongoDB docker image (recommended), follow this next steps:
- Create the private network qa-net running the command `docker network create -d bridge qa-net`.
- Add the MongoDb container to the network running the command `docker run -dit --rm --name mongoDB -p 27017:27017 --network qa-net mongo`, for reference see the docker documentation [a link](https://docs.docker.com/engine/reference/commandline/network_connect)
- Check if the image is runnig on the private network with the command `docker network inspect qa-net`

And thats it, you are ready to go :)



