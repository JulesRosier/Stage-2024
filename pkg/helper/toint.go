package helper

import "strconv"

func StringToInt(s string) int32 {
	i, err := strconv.Atoi(s)
	if err != nil {
		DieMsg(err, "unable to convert string to int")
		return -1
	}
	return int32(i)
}
