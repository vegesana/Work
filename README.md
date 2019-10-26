# Work
# TextProcessingTool.pdf
    goRoutines: We create one go routine for string i.e all the data 
        after that string in the file will go to that go routine. 
        During init, if we don't have any match, then we created one 
        Junk Channel and One go routine listening on it

 
    inputfile.txt 
        This file is read by main go routine. 
    Gch: 
        This variabel is set to the channel on which the data
    strslice: This contains the checkpoints (as string match) which
            will be read from file and will be sent to that specific 
            go routine [goRoutine] for processing.
   Processing in goroutine
       this is regular expression processing and will store the data
            into respective file define dby Data struct
    lineData struct : This contains both teh hostname and the line to 
        get processed.
    mac.go : This file have processing of mac specifi data inside the input
        file
    build : This shell script contians all the input go file to generate 
        the exec
    debug.log : All the debug logs are writted to this file
