# Health Check Service and Mail-Server

This repository contains two core components for managing and maintaining a `mail-server` application:

1. A **Bash-based health check service** that monitors the server’s health, attempts to restart it if unresponsive, and sends email notifications upon repeated failures.
2. A **Go-based server** that handles contact form submissions, saves messages to a database, and sends emails using a local Postfix setup.

These tools work together to ensure high availability and seamless operation of your mail server.

## Repository

**GitHub Repository:** [mail-server](https://github.com/longsizhuo/mail-server)

## Features

### Health Check Service

- Periodic health checks for the `mail-server`, ensuring continuous monitoring of service availability.
- Automatically restarts the `mail-server` if it becomes unresponsive, reducing downtime without manual intervention.
- Sends email alerts if the service fails to recover after a threshold number of attempts, enabling quick response from administrators.
- Automatically frees port 8181 if it is occupied by a rogue process, ensuring the service can bind to the required port.
- Provides detailed logging for troubleshooting and monitoring the status of the service over time.

### Mail-Server Application

- Handles `POST /contact` requests to accept form submissions containing user details.
- Saves submitted messages to a MySQL database using GORM.
- Sends email notifications with form details to a configured recipient using the local `mail` command.
- Includes a `GET /health` endpoint for health checks, returning a simple status response.
- Implements CORS for secure cross-origin resource sharing.

## Prerequisites

- A Linux-based operating system to host the service and application.
- **For Health Check Script:**
  - Bash shell for running the script and automation tasks.
  - `curl` for performing health checks on the configured URL endpoint.
  - `mail` command (e.g., part of `mailutils` or `postfix`) for sending email notifications to administrators.
  - Systemd for service management (optional but highly recommended).
- **For Mail-Server Application:**
  - Go environment for building and running the application.
  - MySQL database for storing contact form submissions.

## Installation

### Mail-Server Application

1. Clone the repository from GitHub:

   ```bash
   git clone https://github.com/longsizhuo/mail-server.git
   cd mail-server
   ```

2. Navigate to the `mail-server` source directory and build the application:

   ```bash
   go build -o mail-server
   ```

3. Deploy the `mail-server` executable to a suitable location:

   ```bash
   sudo cp mail-server /usr/local/bin/mail-server
   sudo chmod +x /usr/local/bin/mail-server
   ```

4. Configure the MySQL database connection in the `server.go` file or through environment variables. Ensure the database is set up and migrations are applied automatically.

5. Start the server manually or configure it as a systemd service (see below).

### Health Check Script

1. Navigate to the health check script directory:

   ```bash
   cd health_check
   ```

2. Make the script executable:

   ```bash
   chmod +x health_check.sh
   ```

3. Edit the script `health_check.sh` to configure the following variables:
   - `URL`: The health check endpoint, which defaults to `http://localhost:8181/health`.
   - `CHECK_INTERVAL`: Interval between health checks in seconds, defaulting to `10` seconds.
   - `FAIL_THRESHOLD`: The maximum number of consecutive failures allowed before action is taken, defaulting to `10` failures.
   - `EMAIL_TO`: The email address where alerts will be sent in case of service failures.
   - `EMAIL_SUBJECT` and `EMAIL_BODY`: Customizable subject and body of the alert email.

4. Deploy the script as a systemd service or run it manually (see below).

## Usage

### Running the Mail-Server Manually

```bash
mail-server
```

This starts the server on `0.0.0.0:8181`. You can customize the port or other configurations as needed.

### Running the Health Check Script Manually

```bash
./health_check.sh
```

This will start the monitoring loop and immediately begin checking the service health based on the configured parameters. The logs will display in the terminal for real-time feedback.

### Setting Up as Systemd Services

#### Mail-Server Service

1. Create a `systemd` service file:

   ```bash
   sudo nano /etc/systemd/system/mail-server.service
   ```

2. Add the following content:

   ```ini
   [Unit]
   Description=Mail-Server Application
   After=network.target

   [Service]
   ExecStart=/usr/local/bin/mail-server
   Restart=always
   User=longsizhuo
   Environment=PATH=/usr/bin:/usr/local/bin

   [Install]
   WantedBy=multi-user.target
   ```

3. Reload systemd and start the service:

   ```bash
   sudo systemctl daemon-reload
   sudo systemctl start mail-server
   sudo systemctl enable mail-server
   ```

#### Health Check Service

1. Create a `systemd` service file:

   ```bash
   sudo nano /etc/systemd/system/health_check.service
   ```

2. Add the following content:

   ```ini
   [Unit]
   Description=Health Check Service for Mail-Server
   After=network.target

   [Service]
   ExecStart=/bin/bash /path/to/health_check.sh
   Restart=always
   User=longsizhuo
   Environment=PATH=/usr/bin:/usr/local/bin

   [Install]
   WantedBy=multi-user.target
   ```

3. Reload systemd and start the service:

   ```bash
   sudo systemctl daemon-reload
   sudo systemctl start health_check
   sudo systemctl enable health_check
   ```

## Troubleshooting

- For `mail-server`, ensure the database connection is configured correctly and accessible. Check logs for any errors during startup or operation.
- For the health check script, ensure the `curl` and `mail` commands are available and properly configured. Review logs to debug failures.
- Use the following commands to check logs for either service:

  ```bash
  sudo journalctl -u mail-server
  sudo journalctl -u health_check
  ```

## License

This project is licensed under the MIT License. See the `LICENSE` file included in the repository for detailed terms and conditions.

## Contributing

We welcome contributions to improve this project! If you encounter any issues or have suggestions for new features, please open an issue or submit a pull request. Contributions to enhance the script, expand documentation, or improve server functionality are greatly appreciated.

---

Ensure any service-specific paths or settings are appropriately updated for your environment.

