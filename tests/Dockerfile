FROM golang:latest 

RUN mkdir /app                        
ADD ./src /app/                        
ADD ./tests /app/                        
WORKDIR /app                         

RUN go get -v "github.com/gorilla/websocket"
RUN go test -v -run ws_test.go                       
