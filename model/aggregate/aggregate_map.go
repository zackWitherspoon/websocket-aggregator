package aggregate

type Duration struct {
	StartTime int64
	EndTime   int64
}

func (dur *Duration) Between(timestamp int64) bool {
	if timestamp > dur.StartTime && timestamp <= dur.EndTime {
		return true
	}
	return false
}
