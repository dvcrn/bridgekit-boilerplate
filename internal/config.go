package internal

import (
	_ "embed"

	"github.com/dvcrn/matrix-bridgekit/bridgekit"
	"go.mau.fi/util/configupgrade"
	"maunium.net/go/mautrix/bridge/bridgeconfig"
)

//go:embed example-config.yaml
var ExampleConfig string

var _ bridgekit.ConfigGetter = &Config{}
var _ bridgeconfig.BridgeConfig = &MyBridgeConfig{}

type MyBridgeConfig struct {
	SomeKey            string                           `yaml:"some_key"`
	Encryption         bridgeconfig.EncryptionConfig    `yaml:"encryption"`
	CommandPrefix      string                           `yaml:"command_prefix"`
	ManagementRoomText bridgeconfig.ManagementRoomTexts `yaml:"management_room_text"`
	DoublePuppetConfig bridgeconfig.DoublePuppetConfig  `yaml:",inline"`
}

func (m MyBridgeConfig) FormatUsername(username string) string {
	//TODO implement me
	return username
}

func (m MyBridgeConfig) GetEncryptionConfig() bridgeconfig.EncryptionConfig {
	//TODO implement me
	return m.Encryption
}

func (m MyBridgeConfig) GetCommandPrefix() string {
	//TODO implement me
	return m.CommandPrefix
}

func (m MyBridgeConfig) GetManagementRoomTexts() bridgeconfig.ManagementRoomTexts {
	//TODO implement me
	return m.ManagementRoomText
}

func (m MyBridgeConfig) GetDoublePuppetConfig() bridgeconfig.DoublePuppetConfig {
	//TODO implement me
	return m.DoublePuppetConfig
}

func (m MyBridgeConfig) GetResendBridgeInfo() bool {
	//TODO implement me
	return false
}

func (m MyBridgeConfig) EnableMessageStatusEvents() bool {
	//TODO implement me
	return false
}

func (m MyBridgeConfig) EnableMessageErrorNotices() bool {
	return true
}

func (m MyBridgeConfig) Validate() error {
	return nil
}

type Config struct {
	*bridgeconfig.BaseConfig `yaml:",inline"`
	BridgeConfig             MyBridgeConfig `yaml:"bridge"`
	AnotherKey               string         `yaml:"another_key"`
}

func (m Config) Base() *bridgeconfig.BaseConfig {
	return m.BaseConfig
}

func (m Config) Bridge() bridgeconfig.BridgeConfig {
	//TODO implement me
	return m.BridgeConfig
}

func (m *Config) DoUpgrade(helper *configupgrade.Helper) {
	bridgeconfig.Upgrader.DoUpgrade(helper)
}
