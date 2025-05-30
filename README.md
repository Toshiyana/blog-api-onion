# Blog Application API with Onion Architecture

## Overview

This is a blog application API designed using the Onion Architecture.

## Architecture

1. **Domain Layer**

   * Contains the core business logic of the application
   * Composed of entities, value objects, domain services, and repositories
   * Independent of external layers
   * Repository

     * Responsible for persisting domain objects
     * Defines only interfaces

2. **Usecase Layer**

   * Contains the business logic to implement the application's features
   * Depends on the domain layer and repository interfaces

3. **Infrastructure Layer**

   * Contains technical implementations such as databases and external APIs
   * Provides implementations for repository interfaces

4. **UI Layer**

   * Provides the user interface
   * Includes HTTP handlers, middleware, etc.

## API Endpoints

### User-related

* `POST /api/users/register` - Register a new user
* `POST /api/users/login` - User login
* `GET /api/users/:id` - Get user information
* `PUT /api/users/:id` - Update user information
* `DELETE /api/users/:id` - Delete user

### Blog-related

* `POST /api/blogs` - Create a new blog post
* `GET /api/blogs` - Get list of blog posts
* `GET /api/blogs/:id` - Get blog post details
* `GET /api/users/:id/blogs` - Get list of blog posts by user
* `PUT /api/blogs/:id` - Update blog post
* `DELETE /api/blogs/:id` - Delete blog post

### Comment-related

* `POST /api/blogs/:id/comments` - Add a comment to a blog post
* `GET /api/blogs/:id/comments` - Get comments for a blog post
* `PUT /api/comments/:id` - Update a comment
* `DELETE /api/comments/:id` - Delete a comment
