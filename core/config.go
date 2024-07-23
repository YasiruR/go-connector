package core

type Config struct {
	Servers struct {
		DSP struct {
			HTTP struct {
				Port int `json:"port"`
			} `json:"http"`
		} `json:"dsp"`
		Gateway struct {
			HTTP struct {
				Port int `json:"port"`
			} `json:"http"`
		} `json:"gateway"`
	} `json:"servers"`
}
