# Blog REST API w/Authorization Using Cerbos

## Technologies Used

- Go
- [Echo Web Framework](https://echo.labstack.com/docs)
- [Cerbos](https://www.cerbos.dev/)

## Overview

This project implements a RESTful API for a blog application with built-in authorization using Cerbos. The API allows users to perform **CRUD (Create, Read, Update, Delete)** operations on blog posts while enforcing access control policies defined by Cerbos.

## Features

- **User Authentication**: Users can authenticate using basic authentication to access protected endpoints.
- **Post Management**: Users can create, read, update, and delete blog posts via the API.
- **Fine-Grained Authorization**: Access to API endpoints is controlled by Cerbos policies, enabling fine-grained control over who can perform specific actions on blog posts based on user roles, resource attributes, and environmental context.

## User Database and Roles

| Users | Passwords | Roles |
| ----- | --------- | ----- |
| kunal | kunalPass | admin |
| bella | bellaPass | user  |
| john  | johnPass  | user  |


## Installation

To run the project locally, follow these steps:

1. Clone the repository:
    ```bash
    git clone https://github.com/verma-kunal/blogapi-auth-cerbos.git
    ```
2. Install dependencies:
    ```bash
    go mod tidy
    ```
3. Build and run the project:
    ```bash
    cerbos run --set=storage.disk.directory=cerbos/policies -- go run main.go
    ```

## Usage

Once the project is running, you can interact with the API using tools like cURL or Postman. Here are some example API endpoints:

- **Create a Post:**
    ```bash
    curl -i -u kunal:kunalPass -X PUT http://localhost:8080/posts -d '{"title": "gitops 101", "owner": "kunal"}'
    ```
- **Read a Post:**
    ```bash
    curl -i -u kunal:kunalPass -X GET http://localhost:8080/posts/1
    ```
- **Update a Post:**
    ```bash
    curl -i -u kunal:kunalPass -X POST http://localhost:8080/posts/1 -d '{"title": "kubernetes 101", "owner": "kunal"}'
    ```
- **Delete a Post:**  
    ```bash
    curl -i -u kunal:kunalPass -X DELETE http://localhost:8080/posts/1
    ```
    
## Contributing

**Contributions are welcome!** If you'd like to contribute to the project, please fork the repository, make your changes, and submit a pull request. Make sure to follow the contribution guidelines outlined in the repository.

## License

This project is licensed under the [Apache License](https://github.com/verma-kunal/blogapi-auth-cerbos/blob/main/LICENSE). See the LICENSE file for details.
    