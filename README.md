# 🚀 GRPC CLI Chat Backend

Welcome to the **GRPC CLI Chat Backend**! This monorepo houses a powerful chat application comprised of three main services: `auth`, `chat`, and a `cli` client. With a robust authentication system using access and refresh tokens, and independent databases for each service, this project is built to scale and impress. Let's dive into the details! 🌟

## Features 🎉

- **Monorepo Architecture**: Three cohesive applications in one repository.
  - `auth`: Handles user authentication and authorization.
  - `chat`: Manages chat messages and chat rooms.
  - `cli`: Command-line interface for interacting with the chat service.
- **Secure Authentication**: Utilizes access and refresh tokens for secure communication.
- **Database Independence**:
  - PostgreSQL for the authentication service.
  - MongoDB for the chat service.
- **Local file for Token Management**: Stores user tokens securely.
- **Comprehensive Makefile**: Simplifies setup, testing, and deployment.

## Getting Started 🚀

### Prerequisites

Ensure you have the following installed on your system:

- Docker & Docker Compose
- Go (latest version)
- Make

### Setup Instructions

1. **Clone the Repository**

```bash
git clone https://github.com/javascriptizer1/grpc-cli-chat.backend.git
cd grpc-cli-chat.backend
```

2. **Environment Configuration**

```bash
cp .env.example .env
```

3. **Start Services with Docker Compose**

```bash
docker-compose up -d
```

4. **Download proto dependencies**

```bash
make vendor-proto
```

5. **Install other cli dependencies**

```bash
make install-deps
```

## Command Usage 🚀

### Auth Service Commands

- **Register a New User**

```bash
go run main.go register --name "John Doe" --email "john@example.com" --password "password" --password-confirm "password"
```

- **Login**

```bash
go run main.go login --login "john@example.com" --password "password"
```

### Chat Service Commands

- **Create a Chat**

```bash
go run main.go create-chat --user-ids="<uuid>"
```

- **List of Chats**

```bash
go run main.go list-chat
```

- **Connect to a Chat**

```bash
go run main.go connect-chat --chat-id "chat123"
```

### CLI Commands Overview

Our CLI client supports various commands to interact with the backend services seamlessly. Here are a few:

- `register`: Register a new user.
- `login`: Login with existing credentials.
- `create-chat`: Create a chat with the specified users.
- `list-chat`: List of all user chats.
- `connect-chat`: Connect to a specific chat by chat ID.

## Architecture Overview 🏗️

Our project is organized as a monorepository containing three main applications:

### Auth Service

- **Database**: PostgreSQL
- **Functionality**: manages user registration, login, and authentication using access and refresh tokens.

### Chat Service

- **Database**: MongoDB
- **Functionality**: handles chat functionalities, including connecting to a chat and sending messages

### CLI Client

- **Framework**: Cobra
- **Functionality**: a command-line interface that allows users to interact with the Auth and Chat services.

#### Key Features

- **Access and Refresh Tokens**: Secure authentication with token-based access control.

- **PostgreSQL and MongoDB**: Two separate databases for different services to ensure scalability and maintainability.

- **Service Provider Pattern**: A clean and modular approach to managing dependencies and service interactions.

- **Interceptors**: Automatically handle token refresh to ensure seamless user experience.

### Development 💻

- **Generate Go from Proto**

```bash
make generate-api
```

- **Up migrations**

```bash
make migrations-up
```

- **Generate migration**

```bash
name=migration_name make migrations-generate
```

- **Lint the Code**

```bash
make lint
```

- **Build the Project**

```bash
make build
```

## Acknowledgments 🙌

I hope you enjoy using the CLI Chat as much as we enjoyed building it. If you find it useful, give us a ⭐ on GitHub!

Happy chatting! 🎉🚀
