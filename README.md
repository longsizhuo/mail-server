# Health Check Service for Mail-Server

This repository contains a Bash-based health check service for monitoring the `mail-server` application. The script periodically checks the health status of the server and attempts to restart the service if it becomes unresponsive. Additionally, it sends email notifications in case of repeated failures. This tool is particularly helpful for maintaining high availability and automating the monitoring process for critical services.

## Features

- Periodic health checks for the `mail-server`, ensuring continuous monitoring of service availability.
- Automatically restarts the `mail-server` if it becomes unresponsive, reducing downtime without manual intervention.
- Sends email alerts if the service fails to recover after a threshold number of attempts, enabling quick response from administrators.
- Automatically frees port 8181 if it is occupied by a rogue process, ensuring the service can bind to the required port.
- Provides detailed logging for troubleshooting and monitoring the status of the service over time.

## Prerequisites

- A Linux-based operating system to host the service.
- Bash shell for running the script and automation tasks.
- `curl` for performing health checks on the configured URL endpoint.
- `mail` command (e.g., part of `mailutils` or `postfix`) for sending email notifications to administrators.
- Systemd for service management, which is optional but highly recommended for ease of deployment and operation.

## Installation

1. Clone the repository from GitHub to your local environment:

   ```bash
   git clone https://github.com/<your-username>/<repository-name>.git
   cd <repository-name>
   ```

2. Make the script executable to ensure it can run properly:

   ```bash
   chmod +x health_check.sh
   ```

3. Edit the script `health_check.sh` to configure the following variables:
   - `URL`: The health check endpoint, which defaults to `http://localhost:8181/health`. Ensure the URL points to the appropriate health endpoint of your service.
   - `CHECK_INTERVAL`: Interval between health checks in seconds, defaulting to `10` seconds.
   - `FAIL_THRESHOLD`: The maximum number of consecutive failures allowed before action is taken, defaulting to `10` failures.
   - `EMAIL_TO`: The email address where alerts will be sent in case of service failures.
   - `EMAIL_SUBJECT` and `EMAIL_BODY`: Customizable subject and body of the alert email to provide detailed information.

4. Deploy the `mail-server` executable and ensure it is running on port 8181. Confirm the service is accessible and responding to health check requests.

## Usage

### Running the Script Manually

To run the health check script manually, use the following command:

```bash
./health_check.sh
```

This will start the monitoring loop and immediately begin checking the service health based on the configured parameters. The logs will display in the terminal for real-time feedback.

### Setting Up as a Systemd Service

1. Create a `systemd` service file to automate the script:

   ```bash
   sudo nano /etc/systemd/system/health_check.service
   ```

2. Add the following content to the service file:

   ```ini
   [Unit]
   Description=Health Check Service for Mail-Server
   After=network.target

   [Service]
   ExecStart=/bin/bash /path/to/health_check.sh
   Restart=always
   User=<your-username>
   Environment=PATH=/usr/bin:/usr/local/bin

   [Install]
   WantedBy=multi-user.target
   ```

3. Reload the `systemd` configuration to register the new service:

   ```bash
   sudo systemctl daemon-reload
   ```

4. Start the service using `systemd`:

   ```bash
   sudo systemctl start health_check
   ```

5. Enable the service to automatically start on system boot:

   ```bash
   sudo systemctl enable health_check
   ```

6. Check the status of the service to confirm it is running correctly:

   ```bash
   sudo systemctl status health_check
   ```

### Disabling the Script

To disable the health check script and prevent it from running:

- If running as a `systemd` service, use the following commands:
  ```bash
  sudo systemctl stop health_check
  sudo systemctl disable health_check
  ```

- If the script was added to `crontab` or `/etc/rc.local`, manually remove or comment out the corresponding entry to disable it.

## Troubleshooting

- Ensure the `mail-server` application is running and accessible on the configured URL endpoint. Test the URL using `curl` to confirm its availability.
- Verify that the `mail` command is properly configured on the server to send email alerts. Test by sending a manual email using:

  ```bash
  echo "Test email" | mail -s "Test Subject" your-email@example.com
  ```

- Review the logs for the health check service to identify potential issues. Use the following command for `systemd` managed services:

  ```bash
  sudo journalctl -u health_check
  ```

- If the script does not run or behaves unexpectedly, ensure it has executable permissions and that all required dependencies (`curl`, `mail`) are installed.

## License

This project is licensed under the MIT License. See the `LICENSE` file included in the repository for detailed terms and conditions.

## Contributing

We welcome contributions to improve this project! If you encounter any issues or have suggestions for new features, please open an issue or submit a pull request. Contributions to enhance the script or expand documentation are greatly appreciated.

---

**Note:** Replace `<your-username>` and `<repository-name>` with your actual GitHub username and repository name before sharing this README publicly. Ensure any service-specific paths or settings are appropriately updated for your environment.

