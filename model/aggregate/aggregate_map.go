package aggregate

type Duration struct {
	StartTime int64
	EndTime   int64
}

func (dur *Duration) Between(ts int64) bool {
	if ts > dur.StartTime && ts <= dur.EndTime {
		return true
	}
	return false
}
