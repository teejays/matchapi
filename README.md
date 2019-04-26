# Match API

This is a simple exercise to create a sample matchmaking app API. The API server is all written in Golang. It makes use of [`gorrila/mux`](https://github.com/gorilla/mux) for HTTP multiplexing and [`teejays/gofildb`](https://github.com/teejays/gofiledb) (more on this below) for a persistant database.

## Installation

The project has been designed to be setup on a dev environment with relatively ease.

### Prerequisites

In order to use start and test the server locally, you will need a few things installed:
1. Golang: Install from the official website [here](https://golang.org).
2. I think that's it, can't think of anything else.

### Installation
The API can be setup using the following commands:
1. Build the binary; `make build`
2. Start the server: `make run`

### Testing
You can run the unit and integration tests using: `make test`

## Documentation

### Entities & Endpoints
There are two main entities in this server: _users_ and _likes_. The codebase,the two services and the API endpoints are structured around these.

#### **USER**

A user object represents the basic user of the API. It is designed in a modular way so it can be broken down into editable and shareable parts. It is implemented in `service/user/v1`. It has the following API endpoints:

- **POST** `/v1/user`: creates a new user and returns the newly created user object including the randomly generated ID: `curl localhost:8080/v1/user -d '{"FirstName":"Jon","LastName":"Doe", "Email": "jon.doe@email.com", "Gender": 1}'`

- **GET** `/<auth_user_id>/v1/user`: provides user obejct of the authenticated user: `curl localhost:8080/<user_id>/v1/user`

- (TODO) **GET** `/<auth_user_id>/v1/user/<user_id>`: provides a non-simplified version of the user object of user with id `user_id`. Personal non-shareable data is excluded: `curl localhost:8080/<user_id>/v1/user`

#### **LIKE**
Like resource represents the action of a user liking another user. It is implemented in `service/like/v1` and `service/like/v2`. It has the following API endpoints:

- **GET** `/<auth_user_id>/v1/like/incoming`: provides all the likes that the authenticated user has received. The data only includes the `userID`'s of user's that have liked the caller, so the the client is expected to make further calls to the API to get more details on any individual user. This behaviour was changed in version 2 of the same API. Sample request: `curl localhost:8080/<user_id>/v1/like/incoming`. 

- **GET** `/<auth_user_id>/v2/like/incoming`: provides all the likes that the authenticated user has received. The data is much richer and field names are more better as compared to V1. This version of the APi also includes some basic user info so the client doesn't have to make subsequent calls to the API for user details. Sample request: `curl localhost:8080/<user_id>/v2/like/incoming`

- **POST** `/<auth_user_id>/v2/like`: represents a new _like_ action. This functionality was not included in previous version of the like API. Sample request: `curl -X "POST" localhost:8080/{userid}/v2/like -d '{"ReceiverID": 3}'`

### Tech Stack

#### WebFramework (None)
I did not use any web framework for the project. The project folder structure and the relationships between different Go packages is something I decided on looking at the how complicated the code was. The simple nature of the project didn't make it necessary to go for a proper framework.


#### Authentication
The authentication layer for this server hasn't been implemented yet. However, the API is built in a way that that Basic Auth could be incorporated easily without changing the structure of the code. This means that the API handlers have the ability figure out what user is making the request, without really affecting RESTful route patterns. This is done by passing in the `userID` of the 'authenticated' user as a prefix to the standard route e.g. `/v1/resource` becomes `/<userID>/v1/resouce`.


#### GoLang
I decided to write this in Go for one main reason: I love Go. I picked up GoLang while at work about 3 years ago. Since then, I have fallen in love the opinionatedness of the language, the strictly typed nature. There is a Go way of doing things, and that's something I've started to admire.

#### Database
I am using my own [_GoFiledb_](https://github.com/teejays/gofiledb) package for as a database. GoFiledb is in alpha (so no documentation) but a simple and personally tested, minimalistic Go client that lets applications use the filesystem as a database. The client is still in development phase. The main advantage of GoFiledb is that it uses the years of optimization efforts that went into file systems to make reading and serving of data fast. It is very quick to set up (vs. a proper database, which are sometimes an overkill for a simple project). 

_Scalability:_
This is a minimalistic API, developed mostly for fun and experimentation reasons. In order to scale it further, a few decisions probably need to be changed. For example, the local file syetem based data storage should probably be replaced by a proper schemaless DB system.


