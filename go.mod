module github.com/prskr/git-age-keyring-agent

go 1.22

toolchain go1.22.1

require (
	buf.build/gen/go/git-age/agent/connectrpc/go v1.16.0-20240327083355-cbf528090598.1
	buf.build/gen/go/git-age/agent/protocolbuffers/go v1.33.0-20240327083355-cbf528090598.1
	connectrpc.com/connect v1.16.0
	connectrpc.com/grpchealth v1.3.0
	connectrpc.com/grpcreflect v1.2.0
	filippo.io/age v1.2.0
	github.com/99designs/keyring v1.2.2
	github.com/adrg/xdg v0.4.0
	github.com/alecthomas/kong v0.9.0
	github.com/lmittmann/tint v1.0.4
	github.com/whilp/git-urls v1.0.0
	golang.org/x/net v0.22.0
)

require (
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/danieljoos/wincred v1.2.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.6.0 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/term v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
