## simple rpc impl like thrift

## Features

* basic type serialize && deserialize
* one tcp connection multiplexing

## Module

### Transport

1. Define how to read & write basic type via io.ReadWriter
1. In this example, we impl `Int32` and `String`

### Protocol

1. Define how message exchange in one round trip
1. In this example, `message` = `seq` + `type` + `payload`
1. `payload` just treat as `string`

### Processor

1. Define how server process request message and use protocal write response back
1. In this example, we impl a simple function multiply

### Client

1. In this example, use tcp connection as `io.ReadWriter`
1. Once connect, start a goroutine loop read response message
1. Client hold map<seq, req>, so that can do multi request and wait response out of order

### Server

1. Same as client, use tcp connection as `io.ReadWriter`
1. Once listen port, create loop wait for new client connect
1. Once new client connection establish, start a goroutine which loop read message.
1. Once read message goruntine got something, start another goroutine use processor handler the message
1. Obviously the server send response to the client out of order

## Example

```bash
go run server.go
go run client.go
```

We do 20 times request at the same time.
In the client ouput log, we can see request and response.


