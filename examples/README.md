# Examples

You can find two practical examples here, client SMPP cli and server SMPP implementation.

## SMSC Server

You can start server with:

    go run server/main.go -addr localhost:2775 -systemid exampleserver

This will start long running process that will listen for tcp connections on _localhost:2775_ and echo received SMPP messages in upper case.

## ESME Client

Client is capable of interacting with the server. You can use client CLI like this:

    go run client/main.go -server localhost:2775 -dst_addr 11111 -src_addr 22222 -msg "This is the message"

This will connect to the server at _localhost:2775_, do initial binding in tranceiver mode and send submit_sm with provided parameters. Upon finishing it will unbind from the server.
