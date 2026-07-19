# litDB

> A lightweight Redis-inspired in-memory key-value database server written in Go.

litDB is a high-performance, asynchronous in-memory database that implements a subset of Redis features while focusing on simplicity, learning, and performance. It uses Linux `epoll` for scalable non-blocking I/O and supports persistence through an Append-Only File (AOF).

---

## ✨ Features

- 🚀 Asynchronous TCP server powered by Linux `epoll`
- 📦 In-memory hash map based key-value storage
- 💾 Append-Only File (AOF) persistence
- ⏳ Key expiration with passive and active cleanup
- 🔄 Command pipelining support
- 🗑️ Basic eviction policy
- 📡 Redis Serialization Protocol (RESP) compatible communication
- ⚡ Handles up to **20,000 concurrent client connections**

---

## Supported Commands

| Command | Description |
|----------|-------------|
| `SET` | Store a key-value pair |
| `GET` | Retrieve a value |
| `DEL` | Delete a key |
| `EXPIRE` | Set expiration for a key |
| `TTL` | Get remaining lifetime of a key |
