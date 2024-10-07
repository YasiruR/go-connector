package boot

func Start() {
	go servers.DSP.Start()
	go servers.Exchanger.Start()
	servers.Gateway.Start()
}
