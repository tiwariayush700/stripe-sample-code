# Accept a Payment

Build a simple checkout form to collect payment details. Included are some basic
build and run scripts you can use to start up the application.

## Data store

Mongodb: for running that on local; one use docker-compose file or run it locally on your own

~~~
docker-compose up
~~~

## Running the sample

1. Run the server for stripe js

~~~
go run server.go
~~~

2. Run the backend server

~~~
go run cmd/server.go
~~~

2. Go to [http://localhost:4242/checkout.html](http://localhost:4242/checkout.html)