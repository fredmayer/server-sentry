
# Server Sentry

**Server Sentry** is a lightweight Go-based CLI tool designed to manage and monitor servers via SSH. It provides information about Docker containers, PM2 processes, Nginx hosts, and more for future. This tool simplifies server management for growing server parks by centralizing essential checks.

---

## Features

- Connect to multiple servers via SSH using a configuration file.
- Check for running Docker containers.
- Verify PM2 processes (if installed).
- List active Nginx hosts from the server configuration.

---

## Installation

### Prerequisites

- [Go](https://golang.org/doc/install) 1.22 or higher installed.
- Access to your servers with SSH credentials or SSH key files.

### Clone the Repository

```bash
git clone https://github.com/fredmayer/server-sentry.git
cd server-sentry
```

### Build the Application

```bash
go install ./cmd/sentry/sentry.go
```

Or download from release for your platform and install

### Run the Application

You can specify the configuration file using the `-c` flag or the `SENTRY_CONFIG_PATH` environment variable:

```bash
sentry -c config.yml
```

Alternatively:

```bash
export SENTRY_CONFIG_PATH=config.yml
sentry
```

---

## Configuration File Format

The configuration file is a YAML file containing a list of servers. Each server requires the following fields:

```yaml
servers:
  - name: "Server 1"
    host: "192.168.1.10"
    port: 22
    user: "root"
    password: "password"

  - name: "Server 2"
    host: "192.168.1.11"
    port: 22
    user: "admin"
    key: "/path/to/private_key"
```

- **name**: Descriptive name for the server.
- **host**: Server's IP address or hostname.
- **port**: SSH port (default is 22).
- **user**: SSH username.
- **password**: SSH password (optional).
- **key**: Path to the private key file (optional).

---

## Usage

### Basic Command

Run the application with the configuration file:

```bash
sentry -c config.yml
```

### Environment Variable Configuration

If the `-c` flag is not provided, the application will use the `SENTRY_CONFIG_PATH` environment variable:

```bash
export SENTRY_CONFIG_PATH=config.yml
sentry
```

### Output Example

```plaintext
Connecting to Server 1 (192.168.1.10)...
Successfully connected to Server 1
Docker containers running on Server 1:
- app-container
- db-container

PM2 is installed on Server 1. Checking processes...
PM2 processes on Server 1:
- API: Running
- Worker: Running

Nginx is installed on Server 1. Checking active hosts...
Active hosts on Server 1:
- example.com
- api.example.com

Connecting to Server 2 (192.168.1.11)...
Successfully connected to Server 2
No running Docker containers on Server 2
PM2 is not installed on Server 2
Nginx is not installed on Server 2
```

---

## Contributing

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a pull request.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---
