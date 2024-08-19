package boot

type Config struct {
	DataSpace struct {
		ParticipantId string `yaml:"participant_id"`
		AssignerId    string `yaml:"assigner_id"`
	}
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
