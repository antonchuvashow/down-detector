package inspector

type InspectorFactory interface {
	Create(*FactoryRegistry, Config) (Inspector, error)

	Marshal(*FactoryRegistry, Inspector) (SerializedConfig, error)
	Unmarshal(*FactoryRegistry, SerializedConfig) (Inspector, error)
}
