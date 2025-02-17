package main

type VLCResponse struct {
	APIVersion   int               `json:"apiversion"`
	AudioDelay   int               `json:"audiodelay"`
	AudioFilters map[string]string `json:"audiofilters"`
	CurrentPLID  int               `json:"currentplid"`
	Equalizer    []interface{}     `json:"equalizer"`
	Fullscreen   int               `json:"fullscreen"`
	Information  struct {
		Category struct {
			Stream0 struct {
				BitsPerSample   string `json:"Bits_per_sample"`
				Channels        string `json:"Channels"`
				Codec           string `json:"Codec"`
				SampleRate      string `json:"Sample_rate"`
				TrackReplayGain string `json:"Track_replay_gain"`
				Type            string `json:"Type"`
			} `json:"Stream 0"`
			Meta struct {
				Title  string `json:"title"`
				Artist string `json:"artist"`
			} `json:"meta"`
		} `json:"category"`
		Chapter int `json:"chapter"`
	} `json:"information"`
	Length   int     `json:"length"`
	Loop     bool    `json:"loop"`
	Position float64 `json:"position"`
	Random   bool    `json:"random"`
	Rate     float64 `json:"rate"`
	Repeat   bool    `json:"repeat"`
	SeekSec  int     `json:"seek_sec"`
	State    string  `json:"state"`
	Stats    struct {
		AveragedemuxBitrate float64 `json:"averagedemuxbitrate"`
		AverageInputBitrate float64 `json:"averageinputbitrate"`
		DecodedAudio        int     `json:"decodedaudio"`
		DecodedVideo        int     `json:"decodedvideo"`
		DemuxBitrate        float64 `json:"demuxbitrate"`
		DemuxCorrupted      int     `json:"demuxcorrupted"`
		DemuxDiscontinuity  int     `json:"demuxdiscontinuity"`
		DemuxReadBytes      int     `json:"demuxreadbytes"`
		DemuxReadPackets    int     `json:"demuxreadpackets"`
		DisplayPictures     int     `json:"displayedpictures"`
		InputBitrate        float64 `json:"inputbitrate"`
		LostABuffers        int     `json:"lostabuffers"`
		LostPictures        int     `json:"lostpictures"`
		PlayedABuffers      int     `json:"playedabuffers"`
		ReadBytes           int     `json:"readbytes"`
		ReadPackets         int     `json:"readpackets"`
		SendBitrate         float64 `json:"sendbitrate"`
		SentBytes           int     `json:"sentbytes"`
		SentPackets         int     `json:"sentpackets"`
	} `json:"stats"`
	SubtitleDelay int    `json:"subtitledelay"`
	Time          int    `json:"time"`
	Version       string `json:"version"`
	VideoEffects  struct {
		Brightness float64 `json:"brightness"`
		Contrast   float64 `json:"contrast"`
		Gamma      float64 `json:"gamma"`
		Hue        float64 `json:"hue"`
		Saturation float64 `json:"saturation"`
	} `json:"videoeffects"`
	Volume int `json:"volume"`
}
