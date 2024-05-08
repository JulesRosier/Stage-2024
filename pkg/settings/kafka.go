package settings

type Kafka struct {
	Brokers       []string `yaml:"brokers"`
	SchemaRgistry struct {
		Urls []string `yaml:"urls"`
	} `yaml:"schemaRegistry"`
	Auth struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"auth"`
	TopicGrouping bool `yaml:"topicGrouping"`
}

func (k *Kafka) SetDefault() {
	k.TopicGrouping = false
}
