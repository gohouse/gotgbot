package bot

type MediaType uint

const (
	MTnil MediaType = iota
	MTphoto
	MTvideo
	MTanimation
	MTaudio
	MTsticker
	MTdocument
)

func (mt MediaType) String() string {
	switch mt {
	case MTphoto:
		return "photo"
	case MTvideo:
		return "video"
	case MTanimation:
		return "animation"
	default:
		return ""
	}

}