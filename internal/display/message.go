package display

type StatusTextMsg struct {
	Text string
}

type BufferSavedMsg struct {
	FileName     string
	BytesWritten int64
}
