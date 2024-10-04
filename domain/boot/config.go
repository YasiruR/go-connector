package boot

type Config struct {
	DataSpace struct {
		ParticipantId string `yaml:"participant_id"`
		AssignerId    string `yaml:"assigner_id"`
		AssigneeId    string `yaml:"assignee_id"`
	} `yaml:"data_space"`
	Catalog struct {
		Title          string   `yaml:"title"`
		Keywords       []string `yaml:"keywords"`
		AccessServices []string `yaml:"access_services"`
		Descriptions   []string `yaml:"descriptions"`
	}
	Servers struct {
		IP  string
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
	DataSources []struct {
		Name     string `yaml:"name"`
		Hostname string `yaml:"hostname"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"data_sources"`
}
