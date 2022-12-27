package resizer

type ResizerInterface interface {
	ResizeImage(name string, qualityIndex float32)
}

type Resize struct{}
