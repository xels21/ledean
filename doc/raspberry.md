# Raspberry

## Service

`/etc/systemd/system/ledean.service`
```ini
[Unit]
Description=LEDean service

[Service]
Type=simple
ExecStart=/home/dean/ledean/start.sh 2>&1 | tee /home/dean/ledean/log.txt
Restart=on-failure
RestartSec=2


[Install]
WantedBy=multi-user.target
```

`sudo systemctl enable ledean.service`
`sudo systemctl start ledean.service`