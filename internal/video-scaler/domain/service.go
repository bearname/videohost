package domain

type VideoScaleService interface {
	PrepareToStream(videoId string, inputVideoPath string, allNeededQualities []Quality, ownerId string) bool
}
