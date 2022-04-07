package aggregate

type Duration struct {
	StartTime int64
	EndTime   int64
}

func (aggDuration *Duration) Between(timestamp int64) bool {
	if timestamp > aggDuration.StartTime && timestamp <= aggDuration.EndTime {
		return true
	}
	return false
}
