# LEDean



## Components
### Webserver:
- proived static file access
- provides REST api

## deps
```
go get github.com/lucasb-eyer/go-colorful
go get github.com/stianeikeland/go-rpio
```

## tools
```
go install github.com/cortesi/modd
go install github.com/hirokidaichi/goviz
```

## Raspberry
build backend with:
>./build_linux.bat
and the frontend with
>ng build

save the output on the raspberry on any folder, eg:

/home/pi/ledean

├── db
│   ├── modeController
│   │   └── index.json
│   ├── ModeRunningLed
│   │   └── parameter.json
│   ├── ModeSolid
│   │   └── parameter.json
│   ├── ModeSolidRainbow
│   │   └── parameter.json
│   └── ModeTransitionRainbow
│       └── parameter.json
├── debug.sh
├── frontend
│   ├── favicon.ico
│   ├── index.html
│   ├── main.js
│   ├── main.js.map
│   ├── MaterialIcons-Regular.eot
│   ├── MaterialIcons-Regular.ttf
│   ├── MaterialIcons-Regular.woff
│   ├── MaterialIcons-Regular.woff2
│   ├── polyfills.js
│   ├── polyfills.js.map
│   ├── runtime.js
│   ├── runtime.js.map
│   ├── scripts.js
│   ├── scripts.js.map
│   ├── styles.js
│   ├── styles.js.map
│   ├── vendor.js
│   └── vendor.js.map
├── kill.sh
├── ledean
├── log.txt
├── start copy.sh
└── start.sh

### Autostart
for Autostart I have used crontab (with root)
sudo crontab -e
>@reboot /home/pi/ledean/start.sh 2>&1 | tee /home/pi/ledean/log.txt

