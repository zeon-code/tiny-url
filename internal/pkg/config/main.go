package config

type Configuration interface {
	Log() Log
	Cache() DatabaseConfiguration
	PrimaryDatabase() DatabaseConfiguration
	ReplicaDatabase() DatabaseConfiguration
	Metric() MetricConfiguration
}

type AppConfiguration struct{}

func NewConfiguration() Configuration {
	return AppConfiguration{}
}

func (c AppConfiguration) Cache() DatabaseConfiguration {
	return RedisConfig{}
}

func (c AppConfiguration) PrimaryDatabase() DatabaseConfiguration {
	return NewPostgresConfig("DB")
}

func (c AppConfiguration) ReplicaDatabase() DatabaseConfiguration {
	return NewPostgresConfig("DB_REPLICA")
}

func (c AppConfiguration) Metric() MetricConfiguration {
	return NewDatadogConfiguration()
}

func (c AppConfiguration) Log() Log {
	return newLogConfig()
}
