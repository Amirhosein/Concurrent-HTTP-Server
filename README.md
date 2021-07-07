# Concurrent HTTP Server

# Introduction

You are to use net/http package and concurrency in golang to implement a simple http server.
First, You are to write a program that reading from specific file and writing into
another file concurrently.
Then, You should implement an api to get file or send file to the user.
REST APIs that you should implement are listed below:

### `localhost:8080/uploadFile`
    
*input formats* :
    
1. json :
   ```json
    {
        "file_id" : integer or string,
        "file" : string
    }
   ```
   
       file_id is an unique 32 bit number and 
       it can be either integer or string.
       It is the name of the file that you should write to.
            
            

# Roadmap

- [ ] Implement `messenger.Dialogs` interface and pass all tests
- [ ] Provide gRPC API for the broker and main functionalities
- [ ] Add basic logs and prometheus metrics
- [ ] Create *dockerfile* and *docker-compose* files for your deployment
- [ ] Deploy your app with the previous `docker-compose` on a remote machine
- [ ] Deploy your app on K8

