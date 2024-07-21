module ledean

go 1.22.5

require (
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket v1.5.1
	github.com/rs/cors v1.10.1
	github.com/sdomino/scribble v0.0.0-20230717151034-b95d4df19aa8
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.1
	periph.io/x/conn/v3 v3.7.0
	periph.io/x/devices/v3 v3.7.0
	periph.io/x/host/v3 v3.8.0
	tinygo.org/x/drivers v0.26.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jcelliott/lumber v0.0.0-20160324203708-dd349441af25 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace tinygo.org/x/drivers => github.com/xels21/tinygo-drivers v0.0.0-20240721163709-8002ea880704
