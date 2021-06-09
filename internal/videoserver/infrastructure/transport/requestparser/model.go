package requestparser

type LikeVideoRequest struct {
	VideoId string
	OwnerId string
	IsLike  bool
}
