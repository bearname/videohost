class VideoStatus {
    static intToStatus(status) {
        switch (parseInt(status)) {
            case 1:
                return "not processed";
            case 2:
                return "processing"
            case 3:
                return "processed"
            default:
                throw new Error("Unknown status");
        }
    }
}

export default VideoStatus