
# HTTP Server With Authentication

# Introduction

You are to use gin framework package and concurrency in golang and jwt-go to implement a simple http server.
You should implement an rest api with authentication on some endpoints to upload or download a file.
All registered users can upload and download files. no user can download other user's files. users can 
give access to other specific users (by their username) to download their files.
REST APIs that you should implement are listed below:

## `localhost:8080/register`
1. **json** :
   ```json
   {
      "username" : "string",
      "password" : "string"
   }
   ```
   create new user
## `localhost:8080/login`
1. **json** :
   ```json
   {
      "username" : "string",
      "password" : "string"
   }
   ```
   logins and gives access token

## `localhost:8080/uploadFile` with auth
uploaded file can be accessed only by uploader in the begining
<br>
***input formats*** :


1. **form** :

         file : []byte

      In this format, `file` is a byte array of the actual file.
      uploads file in filesystem.


***output format*** :

1. **json** :
   
      *successful upload* : 
      ```json
       {
           "download_url" : "string"
       }
      ```

      `download_url` is the path of saved file in filesystem.

      *failure upload* :
      ```json
       {
           "error" : "string"
       }
      ```
      
      `error` is a description of the occurred error.



<br />

## `localhost:8080/downloadFile` with auth


1. **json** :
      ```json
       {
           "download_url" : "string"
       }
      ```

   In this format, `download_url` is an id that we got in successful upload request. it will download the file if user has download premission.


<br />

***output format*** :

1. **json** :

   *successful download* :
      
         http response with actual file
   \
   *failure download* :
      ```json
       {
           "error" : "string"
       }
      ```

   `error` is a description of the occurred error.


## `localhost:8080/addPremission` with auth

1. **json** :
      ```json
       {
           "download_url" : "string",
           "username_to_be_add" : "string"
       }
      ```

   

<br />


# Roadmap
-  Implement middleware for authentication.
-  Set up Postgres with docker
-  Dockerize Server

# hints
- you can handle file permissions with a struct saving allowed user tokens and file details.
- you can implement authentication using jwt with token  in bearer header for requests that needs auth.
- user and file data should be saved in postgres.
