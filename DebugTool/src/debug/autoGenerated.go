package debug

/* This file data is copied after you compile off the shelf
* /Users/rajuv/vgraju/git/Work/Utils - regular expression for counters
* anything chagne
 */

func autoGeneratedCounterExpr() (string, int) {
	counterExp := `\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+)`
	return counterExp, 32
}

func autoGeneratedStatsExpr() string {
	expr := `\s*([\w\s]+):\s*(\d+),([\w\s]+):\s*(\d+),([\w\s]+):\s*(\d+),([\w\s]+):\s*(\d+),([\w\s]+):\s*(\d+),([\w\s]+):\s*(\d+),([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+),([\w\s]+):\s*(\d+),([\w\s]+):\s*(\d+),([\w\s]+):\s*(\d+),([\w\s]+):\s*(\d+),([\w\s]+):\s*(\d+)\s*([\w\s]+):\s*(\d+),([\w\s]+):\s*(\d+),([\w\s]+):\s*(\d+)`
	return expr
}

func autoGeneratedCfgExpr() string {
	expr := `\s*([\w\s]+):\s*(\d+),([\w\s]+):\s*(\d+),([\w\s]+):\s*(\d+),([\w\s]+):\s*(\d+)`
	return expr
}
func autGeneratedNVMEoEUtilExpr() string {
	expr := `\s*(id)\s*:\s*(.*)\s*\s*(type)\s*:\s*(.*)\s*\s*(state)\s*:\s*(.*)\s*\s*(queue)\s*:\s*(.*)\s*\s*(allocated)\s*:\s*(.*)\s*\s*(paired_ctrl_id)\s*:\s*(.*)\s*\s*(local_mac)\s*:\s*(.*)\s*\s*(remote_mac)\s*:\s*(.*)\s*\s*(io_throttle_count)\s*:\s*(.*)\s*\s*(active_io_flag)\s*:\s*(.*)\s*\s*(num_xids)\s*:\s*(.*)\s*\s*(nvmeoe_max_data_length)\s*:\s*(.*)\s*\s*(ctrl shared cookie)\s*:\s*(.*)\s*\s*(ping in progress)\s*:\s*(.*)\s*`
	return expr
}
