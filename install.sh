
#!/bin/bash
set -x
if [ "$(id -u)" != "0" ]; then
  echo "Script needs to be run with sudo" >&2
  exit 1
fi
serviceName="hybrid-sleep"
scriptPath="/usr/local/bin/hybrid-sleep.sh"

# Create the service file
echo "[Unit]
Description=Hybrid sleep service watchdog

[Service]
ExecStart=$scriptPath
Restart=always

[Install]
WantedBy=multi-user.target" | tee /etc/systemd/system/$serviceName.service

cp hybrid-sleep.sh $scriptPath
# Make the script executable
chmod +x $scriptPath

# Reload systemd to recognize the new service
systemctl daemon-reload

# Enable the service to start on boot
systemctl enable $serviceName.service

# Start the service
systemctl start $serviceName.service