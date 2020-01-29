package main

func PclInfoFun(line lineData) {
	Debug("Raju: Postpcl:", line)
	Gdata.GWriteCh <- pcldata{line.filename, line.line}
}
