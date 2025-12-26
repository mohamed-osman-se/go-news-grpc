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
| Validation | `buf.validate`, `protovalidate-go` |
| Streaming | Unary, Server-side, Client-side, Bidirectional |
| Tooling | Buf, Make |
| Storage | In-memory store |




