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
		fmt.Println("Error getting sample rate data:", s)
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
