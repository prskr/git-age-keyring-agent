module github.com/prskr/git-age-keyring-agent

go 1.23

toolchain go1.23.3

require (
	buf.build/gen/go/git-age/agent/connectrpc/go v1.17.0-20240411154421-ccdd2e6e6f4f.1
	buf.build/gen/go/git-age/agent/protocolbuffers/go v1.35.2-20240411154421-ccdd2e6e6f4f.1
	connectrpc.com/connect v1.17.0
	connectrpc.com/grpchealth v1.3.0
	connectrpc.com/grpcreflect v1.2.0
	filippo.io/age v1.2.0
	github.com/99designs/keyring v1.2.2
	github.com/adrg/xdg v0.5.3
	github.com/alecthomas/kong v1.4.0
	github.com/coreos/go-systemd/v22 v22.5.0
	github.com/lmittmann/tint v1.0.5
	github.com/stretchr/testify v1.9.0
	github.com/whilp/git-urls v1.0.0
	golang.org/x/net v0.31.0
)

require (
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/danieljoos/wincred v1.2.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.8.0 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.29.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/term v0.26.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	google.golang.org/protobuf v1.35.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
