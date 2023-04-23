# coinbase-etl

## Instructions for setting up systemd service on your linux host

### First copy of move to the appropriate folder:

$ sudo cp coinbase-etl.service /etc/systemd/system/

### Then reload the systemd daemon to pick up the new service file:

$ sudo systemctl daemon-reload

### Enable your service to start at boot time:

$ sudo systemctl enable coinbase-etl.service

### Start your service:

$ sudo systemctl start coinbase-etl.service

You can now monitor the status of your service with the systemctl status coinbase-etl.service command, and stop it with the systemctl stop coinbase-etl.service command or reload it with sudo systemctl restart coinbase-etl.service
