# Go 1.24.x Utility Project

A reusable Go module with common utilities I frequently use across projects.

## ✨ Features

1. 🔧 Structured logging:
   - Optional HTTP middleware for request logging.
   - Optional HTTP middleware to load single page applications.
     - ⚠️ Warning, do not store sensitive file(s) in SPA directory.
   - Built-in Discord integration for real-time alerts
2. 🧠 In-memory caching:
   - Lightweight, thread-safe, using sync.Map package.
3. ❗Error:
   - Smart error responses based on error type.
   - Utility function for sending standardized HTTP error responses.

## Installation

```code
go get github.com/iTchTheRightSpot/utility
```

### Discord log view

![discord](discord.png)