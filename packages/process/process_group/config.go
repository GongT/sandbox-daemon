package process_group

type NamespaceConfig struct {
	OverlayRoot string `config:"namespace.overlay_root"`
}

func (cfg NamespaceConfig) Validate() error {
	return nil
}
