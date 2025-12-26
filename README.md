<img width="225" height="225" alt="image" src="https://github.com/user-attachments/assets/77e3f999-0e0d-404d-b8c0-df8b3e310b69" />



# News gRPC Service (Go)

This repository contains a **gRPC-based backend service written in Go** that manages news articles using **Protocol Buffers**.  
The project demonstrates how to design, implement, and consume **unary and streaming gRPC APIs** with **schema-level validation** and a **clean Go project structure**.

The goal of this project is to demonstrate **core backend engineering skills**: API design, gRPC communication patterns, validation, and service lifecycle management.

---

## What This Project Does

- Defines a News API using Protocol Buffers (`proto3`)
- Implements a gRPC server in Go
- Implements a gRPC client that exercises all API methods
- Validates requests using rules defined in `.proto` files
- Demonstrates all four gRPC interaction models:
  - Unary RPC
  - Server-side streaming
  - Client-side streaming
  - Bidirectional streaming
- Stores data in an in-memory store for simplicity

---

## gRPC API Overview

### Service: `NewsService`

| Method | RPC Type | Description |
|------|---------|------------|
| `Create` | Unary | Create a news item |
| `Get` | Unary | Retrieve a news item by ID |
| `GetAll` | Server streaming | Stream all news items |
| `UpdateNews` | Client streaming | Update multiple news items in one request |
| `DeletedNews` | Bidirectional streaming | Delete news items via stream |

---

## Technologies Used

| Area | Technology |
|-----|------------|
| Language | Go |
| RPC Framework | gRPC |
| API Definition | Protocol Buffers (proto3) |
| Streaming | Unary, Server-side, Client-side, Bidirectional |
| Tooling | Buf, Make |
| Storage | In-memory store |


## Video Demo

[![Watch the video](https://img.youtube.com/vi/YY7-1twRleU/0.jpg)](https://youtu.be/YY7-1twRleU?si=0WdRf5yrPwzibJ8a)



### How to Run

## Prerequisites

- Go 1.20+
- `make`
- Linux or macOS  
  *(Windows supported via WSL)*

---

### Clone, Setup & Run

Follow these steps to clone the repository, set up dependencies, and run the server and client:

```bash
# Clone the repository
git clone https://github.com/mohamed-osman-se/go-news-grpc.git
cd go-news-grpc

# Install tools, tidy modules, and generate proto code
make setup-run

# Run the gRPC server
go run ./cmd/server/main.go

# In a separate terminal, run the gRPC client
go run ./cmd/client/main.go
```

### Client Demonstrations

The client demonstrates the following operations:

- Validation failures
- Create requests
- Server-side streaming reads
- Client-side streaming updates
- Bidirectional streaming deletes

---

### Health Checks and Shutdown

- The server registers a **gRPC health service**.
- OS signals (`SIGINT`, `SIGTERM`) are handled to allow a **clean shutdown**.
- Ongoing requests are allowed to **complete before exit**.

---

### Limitations

- Data is stored **in-memory** (no persistence).
- No **authentication** or **authorization**.
- No **TLS** (uses insecure credentials for local development).

> These trade-offs are intentional to focus on **gRPC mechanics and API design**.






