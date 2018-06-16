package config

var Config = struct {
	GitHub struct {
		ClientID     string `required:"true" yaml:"client_id"`
		ClientSecret string `required:"true" yaml:"client_secret"`
	}
	Secret string `required:"true"`
}{}
