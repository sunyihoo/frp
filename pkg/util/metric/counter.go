package metric

type Counter interface {
	Count() int32
	Inc(int32)
	Dec(int32)
	Snapshot() Counter
	Clear()
}
