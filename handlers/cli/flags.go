package cli

type ServiceNameFlag struct {
	ServiceName string `env:"GIT_AGE_SERVICE_NAME" name:"service.name" short:"s" default:"git-age"`
}
