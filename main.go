package main

/*
#cgo LDFLAGS: -framework CoreAudio
#include <CoreAudio/CoreAudio.h>
#include <stdio.h>

OSStatus setDefaultOutputDeviceSampleRate(Float64 sampleRate) {
    // Get the default output device
    AudioDeviceID deviceID = kAudioObjectUnknown;
    UInt32 size = sizeof(deviceID);
    AudioObjectPropertyAddress address = {
        kAudioHardwarePropertyDefaultOutputDevice,
        kAudioObjectPropertyScopeOutput,
        kAudioObjectPropertyElementMain
    };

    OSStatus status = AudioObjectGetPropertyData(kAudioObjectSystemObject, &address, 0, NULL, &size, &deviceID);
    if (status != noErr || deviceID == kAudioObjectUnknown) {
		printf("Error getting the default output device");
        return status; // Error getting the default output device
	}

    // Set the sample rate for the default output device
    address.mSelector = kAudioDevicePropertyNominalSampleRate;

    status = AudioObjectSetPropertyData(deviceID, &address, 0, NULL, sizeof(sampleRate), &sampleRate);
    return status;
}
*/
import "C"

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// VLCResponse Define structs to match the JSON response

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

func setSampleRate(sampleRate float64) error {
	// Call the C function to set the sample rate
	status := C.setDefaultOutputDeviceSampleRate(C.Float64(sampleRate))
	if status != 0 {
		return fmt.Errorf("error setting sample rate: %v", status)
	}
	return nil
}

func getStatus() (*VLCResponse, error) {
	url := "http://localhost:8080/requests/status.json"
	userName := ""
	password := "vlcpassword"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	auth := userName + ":" + password
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", basicAuth)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error executing request:", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error executing request:", resp.StatusCode)
		return nil, err
	}
	var result VLCResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response:", err)
		return nil, err
	}
	return &result, nil
}

func convertStringToFloat(s string) float64 {
	ss := strings.Split(s, " ")
	if len(ss) > 0 {
		sampleRate, err := strconv.ParseFloat(ss[0], 64)
		if err != nil {
			fmt.Println("Error converting sample rate to float:", err)
			return float64(48000)
		}
		return sampleRate
	} else {
		fmt.Println("Error converting sample rate to float:", s)
		return float64(48000)
	}

}

func main() {
	runtime.GOMAXPROCS(1)

	var currentSong string
	var currentSampleRate = float64(48000)
	//var currentBitrate float64

	for {
		vlc, err := getStatus()
		if err != nil {
			fmt.Println("Error getting status:", err)
		}

		if currentSong != vlc.Information.Category.Meta.Title {
			fmt.Println("==========================================================================")
			fmt.Println("Song Changed!")
			fmt.Printf("Song Title: %s - Artist: %s\n", vlc.Information.Category.Meta.Title, vlc.Information.Category.Meta.Artist)
			currentSong = vlc.Information.Category.Meta.Title
			newSampleRate := convertStringToFloat(vlc.Information.Category.Stream0.SampleRate)
			if currentSampleRate != newSampleRate {
				currentSampleRate = newSampleRate
				err := setSampleRate(newSampleRate)
				if err != nil {
					fmt.Println("Error setting sample rate:", err)
					currentSampleRate = float64(48000)
				}
				fmt.Printf("Setting sample rate: %.0f Hz\n", newSampleRate)
			} else {
				fmt.Println("Sample rate stayed the same")
			}
			fmt.Println("==========================================================================")
		}
		time.Sleep(2 * time.Second)
	}

}
