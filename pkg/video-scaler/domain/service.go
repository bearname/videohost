package domain

type ScalerService interface {
	PrepareToStream(videoId string, inputVideoPath string, allNeededQualities []Quality, ownerId string) bool
}
