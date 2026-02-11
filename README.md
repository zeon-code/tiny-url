# Tiny URL

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
For more details, see the [doc/architecture](doc/architecture.md) documentation.

---

## 🧪 Getting Started

### 📦 Prerequisites

You’ll need:

- Go (1.24+ recommended)
- The usual Go toolchain

Clone the repo:

```bash
git clone https://github.com/zeon-code/tiny-url.git
```

### 🚀 Running the Service

Navigate to the project directory:

```bash
cd tiny-url
```

Run the service:

```bash
make run
```
The service will start on `http://localhost:8080`.

### 🧪 Testing
Run tests with:

```bash
make test
```

This will execute all unit tests and display the results.