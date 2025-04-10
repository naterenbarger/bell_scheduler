# Setting up Bell Scheduler as a Systemd Service

This guide explains how to set up the Bell Scheduler application to run automatically at startup on your Raspberry Pi using systemd.

## Prerequisites

- Raspberry Pi with Raspberry Pi OS (or compatible Linux distribution)
- Bell Scheduler application installed and built
- Root or sudo access on your Raspberry Pi

## Installation Steps

1. **Build the application**

   Make sure you've built the backend application:

   ```bash
   cd ~/bell_scheduler/backend
   make deps
   make build
   ```

2. **Copy the service file**

   Copy the systemd service file to the systemd directory:

   ```bash
   sudo cp ~/bell_scheduler/bell_scheduler.service /etc/systemd/system/
   ```

3. **Adjust file paths if necessary**

   Edit the service file if your installation directory is different from `/home/pi/bell_scheduler`:

   ```bash
   sudo nano /etc/systemd/system/bell_scheduler.service
   ```

   Update the following lines to match your installation:
   - `WorkingDirectory=/home/pi/bell_scheduler/backend`
   - `ExecStart=/home/pi/bell_scheduler/backend/bin/bell_scheduler`
   - `EnvironmentFile=/home/pi/bell_scheduler/backend/.env`

4. **Reload systemd**

   After copying or editing the service file, reload the systemd daemon:

   ```bash
   sudo systemctl daemon-reload
   ```

5. **Enable the service**

   Enable the service to start automatically at boot:

   ```bash
   sudo systemctl enable bell_scheduler.service
   ```

6. **Start the service**

   Start the service immediately:

   ```bash
   sudo systemctl start bell_scheduler.service
   ```

## Managing the Service

- **Check service status**:
  ```bash
  sudo systemctl status bell_scheduler.service
  ```

- **Stop the service**:
  ```bash
  sudo systemctl stop bell_scheduler.service
  ```

- **Restart the service**:
  ```bash
  sudo systemctl restart bell_scheduler.service
  ```

- **View service logs**:
  ```bash
  sudo journalctl -u bell_scheduler.service
  ```
  
  To follow logs in real-time:
  ```bash
  sudo journalctl -u bell_scheduler.service -f
  ```

## Troubleshooting

- **Service fails to start**:
  - Check the logs: `sudo journalctl -u bell_scheduler.service -n 50`
  - Verify the executable path is correct
  - Ensure the .env file exists and has correct permissions
  - Check that the user specified in the service file has access to the required directories

- **GPIO access issues**:
  - Ensure the user running the service (default: pi) has access to GPIO
  - You may need to add the user to the gpio group: `sudo usermod -a -G gpio pi`

- **Permission issues**:
  - Check file permissions: `ls -la ~/bell_scheduler/backend/bin/`
  - Make the executable file executable: `chmod +x ~/bell_scheduler/backend/bin/bell_scheduler`