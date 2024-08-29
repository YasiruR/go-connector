package boot

func Start() {
	go servers.DSP.Start()
	servers.Gateway.Start()
}
