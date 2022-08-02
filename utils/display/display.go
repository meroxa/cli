package display

func truncateString(oldString string, l int) string {
	str := oldString

	if len(oldString) > l {
		str = oldString[:l] + "..."
	}

	return str
}
