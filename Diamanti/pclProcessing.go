package main

func processPclInfo(line lineData) {
	Error("Raju: Postpcl:", line)
	Gdata.GWriteCh <- pcldata{line.filename, line.line}
}
