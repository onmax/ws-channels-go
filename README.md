# Channels using WebSockets in Go

This project is a backend websocket server. When a client connects to the system, the client could create his own channel or join other channel. The channels in the documentation and in the code are referred to rooms. One member only is in one room at the time. For further documentation, go to: [here](https://onmax.github.io/ws-channels-go/)

The purpose of this project, is to allow me to have a websocket server for further projects that I will make using Flutter.

## Run the application

You can use the ```Makefile``` and run ```make help``` to see all the options. To run the application using Docker, run:

```make docker-run```

It will create an image and run a container with the application serving on port 8080. You can change the port in the file ```.env```

Otherwise, you can run the application if you have golang in your machine using ```go run``` or ```go build```.

## Documentation

You can see the documentation [here](https://onmax.github.io/ws-channels-go/)

## Example

You can go to /example to see an implementation of a websocket client in JS. You will need to run the application in localhost.
