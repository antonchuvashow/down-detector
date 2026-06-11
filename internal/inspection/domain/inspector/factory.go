package inspector

type Factory interface {
	Marshal(*FactoryRegistry, Inspector) (SerializedConfig, error)
	Unmarshal(*FactoryRegistry, SerializedConfig) (Inspector, error)
}
