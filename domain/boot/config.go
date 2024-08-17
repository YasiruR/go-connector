package boot

type Config struct {
	Catalog struct {
		Title          string   `yaml:"title"`
		Keywords       []string `yaml:"keywords"`
		AccessServices []string `yaml:"access_services"`
		Descriptions   []string `yaml:"descriptions"`
	}
	Servers struct {
		DSP struct {
			HTTP struct {
				Port int `yaml:"port"`
			} `yaml:"http"`
		} `yaml:"dsp"`
		Gateway struct {
			HTTP struct {
				Port int `yaml:"port"`
			} `yaml:"http"`
		} `yaml:"gateway"`
	} `yaml:"servers"`
}
