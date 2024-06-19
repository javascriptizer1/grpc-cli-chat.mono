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
- Go (1.22.4)
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

- **List of Users**

```bash
go run main.go list-user
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

## Development 💻

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

## Deployment ⚙️

### Kubernetes Deployment

The application is deployed in a Kubernetes cluster, ensuring high availability, scalability, and efficient resource management. The deployment process is automated using GitHub Actions.

### VPS Deployment

In addition to Kubernetes, the application can also be deployed on a Virtual Private Server (VPS). This flexibility allows you to choose the deployment target that best fits your infrastructure and scaling needs. The deployment process, whether to Kubernetes or a VPS, is managed through GitHub Actions workflows.

### GitHub Actions CI/CD

I use GitHub Actions for continuous integration and deployment. The workflow includes:

1. **Linting**: Code is checked using `golangci-lint`.
2. **Building Docker Images**: Docker images are built for each service.
3. **Pushing Docker Images**: The built images are pushed to Docker Hub.
4. **Deploying to Kubernetes**: The images are deployed to a Kubernetes cluster using Helm charts.
5. **Deploying to VPS**: The application is built, environment variables are configured, and binaries are transferred and set up on the VPS.

### Modern Deployment Practices

Deployment setup exemplifies modern DevOps practices by:

- **Automated CI/CD**: Minimizing manual intervention and reducing the risk of errors.
- **Environment Configurations**: Supporting multiple environments (development, production) with ease.
- **Scalability and Reliability**: Leveraging Kubernetes for managing and scaling microservices effectively.
- **Flexibility**: Offering deployment options to both Kubernetes and VPS, catering to different infrastructure needs.
- **Security**: Storing sensitive information like Docker Hub credentials, Kubernetes config, and VPS SSH keys securely in GitHub Secrets.

### Environment Variables

Ensure that the following environment variables are set in your GitHub repository secrets:

- `DOCKER_HUB_USERNAME`: Your Docker Hub username.
- `DOCKER_HUB_PASSWORD`: Your Docker Hub password.
- `KUBE_CONFIG`: Your Kubernetes configuration file content.
- `NAMESPACE`: Your Kubernetes namespace.
- `VPS_SSH_KEY`: Your VPS SSH private key.
- `VPS_USER`: Your VPS user.
- `VPS_HOST`: Your VPS host address.

## Future Plans 🌅

My plan to rewrite the CLI application from Cobra to Bubble Tea to build a full-fledged console UI. This will provide a more interactive and user-friendly experience for users interacting with my services via the command line.

### Why Bubble Tea?

- **Rich Interactive UI**: Bubble Tea allows us to build rich, interactive UIs directly in the terminal, enhancing user experience.
- **Customization**: Offers a high degree of customization, enabling us to tailor the interface to my specific needs.
- **Active Community**: Supported by an active community and regularly updated. [Bubble Tea GitHub Repository](https://github.com/charmbracelet/bubbletea)
- **Ease of Use**: Simplifies the development of complex terminal applications, making the process more intuitive and efficient.

With these enhancements, I aim to provide a powerful and engaging CLI tool that leverages the capabilities of Bubble Tea to its fullest potential.

## Acknowledgments 🙌

I hope you enjoy using the CLI Chat as much as we enjoyed building it. If you find it useful, give us a ⭐ on GitHub!

Happy chatting! 🎉🚀
