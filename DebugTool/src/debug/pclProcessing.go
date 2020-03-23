package debug

func PclInfoFun(line lineData) {
	Debug("Raju: Postpcl:", line)
	writeToDb(Pcldata{line.Server, line.Line})
}
