# Tiny URL [![Build and test server](https://github.com/zeon-code/tiny-url/actions/workflows/test.yml/badge.svg)](https://github.com/zeon-code/tiny-url/actions/workflows/test.yml) [![CodeQL](https://github.com/zeon-code/tiny-url/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/zeon-code/tiny-url/actions/workflows/github-code-scanning/codeql)

A high‑performance, lightweight URL shortening service written in **Go**.  
Designed with **simplicity, scalability, and extensibility** in mind.

---

## 📖 API Documentation

*   **Interactive Playground:** [Open Documentation](https://zeon-code.github.io/tiny-url/)
*   **OpenAPI Spec:** [`openapi.spec.yaml`](./docs/openapi.spec.yaml)

---

## 🚀 Features

✔️ Shorten long URLs into short, memorable links  
✔️ Fast and efficient implementation in Go  
✔️ Clean, pragmatic layered architecture  
✔️ Designed for extension (analytics, metrics, etc.)

> 🔧 Current core functionality focuses on URL shortening & redirection.

---

## 🧠 Architecture

This project uses a simple layered architecture that balances readability with performance.  
For more details, see the [docs/architecture](doc/architecture.md) documentation.

---

## 📦 Getting Started

####  Prerequisites

You’ll need:

- Go (1.24+ recommended)
- The usual Go toolchain

Clone the repo:

```bash
git clone https://github.com/zeon-code/tiny-url.git
```

#### Running the Service

Navigate to the project directory:

```bash
cd tiny-url
```

Run the service:

```bash
make run
```
The service will start on `http://localhost:8080`.

#### Testing
Run tests with:

```bash
make test
```

This will execute all unit tests and display the results.

#### Migration
Run migrations with:

```bash
make migrate
```

This will execute all migrations and display the results. In order to create a new migration, use the following command:

```bash
make new-migration name=add_users_table
```

This will create a new migration file in the `migrations` directory with the name `add_users_table`.

#### Live documentation (Swagger UI)

Access the interactive API documentation at:

```
http://localhost
```