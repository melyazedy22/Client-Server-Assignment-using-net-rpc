# Assignmnnet 04: Simple Chatroom

Using Demo 1 create a simple chatroom

1. Client

   - the client will dial the rpc of the coordinating server
   - the client will call the remote procedure on the server to send a message
   - the client can fetch all of the messages history from the server using remote procedure call

2. Server

   - every message sent to the server has to be stored in long list
   - the server should return the chat history to the client when the client sends a message

## Notes

- The client should be running forever until you either type the message "exit" or hit Ctrl+c
- `fmt.Scan` reads input token by token, where tokens are separated by spaces. So look for something that will read the whole message if separated by spaces.
- Don't forget to handle the errors if the server goes down while the client is making the connection.
