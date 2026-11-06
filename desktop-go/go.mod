module xcpng-center

go 1.27

require (
	github.com/BurntSushi/toml v1.6.0
	github.com/coder/websocket v1.8.14
	github.com/modelcontextprotocol/go-sdk v1.8.0
	github.com/spf13/cobra v1.10.2
	github.com/wailsapp/wails/v3 v3.1.0
	go/xenapi v0.0.0
	golang.org/x/term v0.47.0
	software.sslmate.com/src/go-pkcs12 v0.7.1
)

require (
	github.com/ebitengine/purego v0.11.1 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/godbus/dbus/v5 v5.2.2 // indirect
	github.com/google/jsonschema-go v0.4.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/segmentio/asm v1.1.3 // indirect
	github.com/segmentio/encoding v0.5.4 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	golang.org/x/crypto v0.57.0 // indirect
	golang.org/x/oauth2 v0.35.0 // indirect
	golang.org/x/sync v0.23.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/time v0.15.0 // indirect
)

replace go/xenapi => ./sdk/xenapi
