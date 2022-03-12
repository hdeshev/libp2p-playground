# Send and Receive Protobuf Messages

Define the messages in a protobuf file, generate code and use it with libp2p transports.

Create an imaginary time service that receives a greeting string and responds with a timestamp.

Run the server with

```
$ go run main.go
```

Copy the server address and run the client with:

```
$ go run main.go -s <server address>
```
