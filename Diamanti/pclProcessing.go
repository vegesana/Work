package main

func PclInfoFun(line lineData) {
	Error("Raju: Postpcl:", line)
	Gdata.GWriteCh <- pcldata{line.filename, line.line}
}
