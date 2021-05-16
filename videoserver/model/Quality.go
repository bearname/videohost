package model

type Quality int

const (
	Q1080p Quality = iota
	Q720p
	Q480p
	Q320p
	Q144p
)

func (q Quality) String() string {
	return [...]string{"1080p", "720p", "480p", "320p", "144p"}[q]
}

func (q Quality) Int() int {
	return [...]int{1080, 720, 480, 320, 144}[q]
}
