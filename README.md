Description
This project implements a configuration management service that handles configuration data using RESTful API endpoints. It supports operations like creating, retrieving, updating, and deleting configuration data and groups. The service is rate-limited and uses a gorilla/mux router for HTTP handling, with Docker containerization provided.

Features
Configuration Management: Manage configurations and config groups.
Rate Limiting: Limits requests to 10 per minute with a burst capacity of 3.
Graceful Shutdown: Supports graceful shutdown on termination signals.
Technologies Used
Go (Golang): Backend language.
Consul: Used for configuration storage.
Docker: Containerization of the application.
Gorilla Mux: Routing for HTTP endpoints.
x/time/rate: Rate limiter for API requests.
