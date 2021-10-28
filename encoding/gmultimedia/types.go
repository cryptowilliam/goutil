package gmultimedia

type MediaFmt int

const (
	MediaFmtUnknown MediaFmt = -1
	MediaFmtPNG     MediaFmt = 0
	MediaFmtSVG     MediaFmt = 1
	MediaFmtJPG     MediaFmt = 2
	MediaFmtWEBP    MediaFmt = 3
	MediaFmtGIF     MediaFmt = 4
	MediaFmtFLV     MediaFmt = 5
	MediaFmtMOV     MediaFmt = 6
	MediaFmtMP4     MediaFmt = 7
)

func DetectFile(filename string) (MediaFmt, error) {
	return MediaFmtUnknown, nil
}

func DetectBuffer(file []byte) (MediaFmt, error) {
	return MediaFmtUnknown, nil
}

var SuffixsOfImage = []string{".png", ".jpeg", ".jpg", ".bmp", ".gif"}
var SuffixsOfVideo = []string{".mp4", ".mp5", ".h264", ".h265", ".flv", ".webm", ".mkv", ".rm", ".rmvb", ".mpg", ".mpeg", ".avi", ".wmv", ".asf", ".mov", ".qt"}
var SuffixsOfAudio = []string{".mp3", ".aac", ".acm", ".aif", ".aifc", ".flac", ".wma", ".wav", ".midi"}
