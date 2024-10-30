package progress

type ProgressFactory interface {
	Progress() (id, name string, weak bool)
}
