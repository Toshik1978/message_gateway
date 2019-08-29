package service

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	configName = "message_gateway.conf"
)

// Vars hold all variables for service running
type Vars struct {
	HTTPAddress string

	ProxyAddress string
	ProxyLogin   string
	ProxyPass    string

	TelegramToken string

	SMSIncomingPath string
	SMSOutgoingPath string
	SMSTempPath     string

	EmailName    string
	EmailSubject string
	EmailSMTP    string
	EmailPort    string
	EmailLogin   string
	EmailPass    string

	UpsAddress string
	UpsName    string
	UpsLogin   string
	UpsPass    string

	PingHosts string
}

// LoadConfig load's configuration file
func LoadConfig(logger *zap.Logger) Vars {
	viper.SetConfigName(configName)
	viper.AddConfigPath("configs")
	viper.AddConfigPath("/etc")
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("Failed to read config file", zap.Error(err))
	}

	return Vars{
		HTTPAddress:     viper.GetString("http.address"),
		ProxyAddress:    viper.GetString("proxy.address"),
		ProxyLogin:      viper.GetString("proxy.login"),
		ProxyPass:       viper.GetString("proxy.pass"),
		TelegramToken:   viper.GetString("telegram.token"),
		SMSIncomingPath: viper.GetString("sms.incoming"),
		SMSOutgoingPath: viper.GetString("sms.outgoing"),
		SMSTempPath:     viper.GetString("sms.temp"),
		EmailName:       viper.GetString("email.from"),
		EmailSubject:    viper.GetString("email.subject"),
		EmailSMTP:       viper.GetString("email.smtp"),
		EmailPort:       viper.GetString("email.smtp_port"),
		EmailLogin:      viper.GetString("email.login"),
		EmailPass:       viper.GetString("email.pass"),
		UpsAddress:      viper.GetString("ups.address"),
		UpsName:         viper.GetString("ups.name"),
		UpsLogin:        viper.GetString("ups.login"),
		UpsPass:         viper.GetString("ups.pass"),
		PingHosts:       viper.GetString("ping.hosts"),
	}
}
