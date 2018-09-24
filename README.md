# go-api-storage-mongodb
API which stores the requests data into mongo db.

##### How to build and run the application:

To run the application you would need docker installed. Install docker first.
Unzip the package

```
$ cd go-api-storage-mongodb/init-database

//Create a docker custom network
$ sudo docker network create --subnet=172.18.0.0/16 mycustomnet

//build a image with dockerfile
$ sudo docker build -t mongo-db .

//Run the container. Which starts the db service on port 27017

$ sudo docker run -d -p 27017:27017 --net mycustomnet --ip 172.18.0.10 mongo-db

$ cd ..

//build a image with dockerfile
$ sudo docker build -t api .

//Run the container. Which starts the API service on port 8080
sudo docker run -d -p 8080:8080 --net mycustomnet --ip 172.18.0.20 api

```


On the host machine open browser and visit http://localhost:8080 .




##### Functionalities:

Assuming CandidateNames are unique.

Please the below commands for appropriate operations.

1) Create Candidate/User
```
    curl -H "Accept: application/json" -H "Content-type: application/json" -X POST -d '{"name":"user1","age":"23","phone":"000-000-0000"}' http://localhost:8080/createUser

    curl -H "Accept: application/json" -H "Content-type: application/json" -X POST -d '{"name":"user2","age":"27","phone":"001-000-1000"}' http://localhost:8080/createUser

    curl -H "Accept: application/json" -H "Content-type: application/json" -X POST -d '{"name":"user3","age":"32","phone":"100-100-1000"}' http://localhost:8080/createUser

```

2) Delete Candidate/User

```
    curl http://localhost:8080/removeUser/?username=user3

```

3) Create or Update Candidate/User Tags(key-value pairs)

```
curl -H "Accept: application/json" -H "Content-type: application/json" -X POST -d '{"lastname":"lastuser1","state":"CA"}' http://localhost:8080/create-updateTags/?username=user1

```
4) Delete Candidate/User Tags(key-value pairs)
```
        curl "http://localhost:8080/deleteTags/?username=user1&deletetags=age,lastname"

```
5) Listall Candidate/User Tags(key-value pairs)

```
          curl http://localhost:8080/listTags/?username=user1

```