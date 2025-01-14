module github.com/vanilla-os/vib

go 1.21

require github.com/spf13/cobra v1.8.0

require github.com/mitchellh/mapstructure v1.5.0

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/vanilla-os/vib/api v0.0.0-20240331150207-852011e4d96f
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/vanilla-os/vib/api => ./api
