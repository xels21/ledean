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



### Config
> The SPI implementation is sensitive to variations in SPI clock speed. On the Raspberry Pi, you will need to add `core_freq=250` to /boot/config.txt to prevent glitching.
> You may also need to increase your SPI buffer size to 12*num_pixels+3, or just max it out with `spidev.bufsize=65536`. That should allopw you to buffer over 5400 Neopixels.
>               - [periph.io - nezled](https://pkg.go.dev/periph.io/x/periph/experimental/devices/nrzled)
### Autostart
for Autostart I have used crontab (with root)
`sudo crontab -e`
>@reboot sleep 10 /home/pi/ledean/start.sh 2>&1 | tee /home/pi/ledean/log.txt

`sleep 10` is used to give the uC 10 seconds time to start SPI, avoiding this error:
`spireg: no port found; did you forget to call Init()?"`

### Pin
for data pin use SPI MOSI (depending on Raspi)
see `doc\raspi_pin_ascii.txt`

for (switch) button, you need to connect it as followed:
```
          _____
gpio_x---|__Ω__|---3.3v

```