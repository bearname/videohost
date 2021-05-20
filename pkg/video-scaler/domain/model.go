package domain

type Quality int

const (
	Q2160p Quality = iota
	Q1440p
	Q1080p
	Q720p
	Q480p
	Q320p
	Q144p
)

func (q Quality) String() string {
	return [...]string{"2160", "1440", "1080", "720", "480", "320", "144"}[q]
}

func (q Quality) Values() int {
	return [...]int{2160, 1440, 1080, 720, 480, 320, 144}[q]
}
