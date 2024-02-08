package engine

import (
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type Voice struct {
	streamer beep.Streamer
	Paused   bool
	Volume   float64
}

func (v Voice) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = v.streamer.Stream(samples)
	if v.Paused {
		for i := range samples {
			samples[i] = [2]float64{}
		}
		return len(samples), true
	} else {
		gain := math.Pow(2, v.Volume)
		for i := range samples[:n] {
			samples[i][0] *= gain
			samples[i][1] *= gain
		}
	}

	return n, ok
}

func (v Voice) Err() error {
	return v.streamer.Err()
}

var soundBank map[string]beep.Buffer

var looped map[string]*Voice

var isInitialised bool = false

func LoadSound(filepath, name string) {
	if _, ok := soundBank[name]; ok {
		log.Println("Sound already loaded: ", name)
		return
	}
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	var streamer beep.StreamSeekCloser
	var format beep.Format

	filetype := strings.Split(filepath, ".")[1]
	if filetype == "mp3" {
		streamer, format, err = mp3.Decode(f)
	} else if filetype == "wav" {
		streamer, format, err = wav.Decode(f)
	}
	if err != nil {
		log.Fatal(err)
	}

	if !isInitialised {
		soundBank = make(map[string]beep.Buffer)
		looped = make(map[string]*Voice)
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		isInitialised = true
	}

	buffer := *beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()
	soundBank[name] = buffer
}

func PlaySound(name string, volume float64) {
	sound, ok := soundBank[name]
	if !ok {
		log.Println("sound does not exist: ", name)
		return
	}

	streamer := sound.Streamer(0, sound.Len())
	voice := Voice{streamer: streamer, Volume: volume, Paused: false}
	speaker.Play(voice)
}

func LoopSound(name string, volume float64) {
	sound, ok := soundBank[name]
	if !ok {
		log.Println("sound does not exist: ", name)
		return
	}

	existing, ok := looped[name]
	if !ok {
		streamer := sound.Streamer(0, sound.Len())
		loop := beep.Loop(-1, streamer)
		looped[name] = &Voice{streamer: loop, Paused: false, Volume: volume}
		speaker.Play(looped[name])
	} else {
		speaker.Lock()
		existing.Paused = false
		speaker.Unlock()
	}
}

func StopLoop(name string) {
	loop, ok := looped[name]
	if !ok {
		log.Println("sound is not actively looping: ", name)
		return
	}
	speaker.Lock()
	loop.Paused = true
	speaker.Unlock()
	delete(looped, name)
}

func PauseLoop(name string) {
	loop, ok := looped[name]
	if !ok {
		log.Println("sound is not actively looping: ", name)
		return
	}

	speaker.Lock()
	loop.Paused = true
	speaker.Unlock()
}

func ClearSounds() {
	speaker.Clear()
}

func close() {
	ClearSounds()
	speaker.Close()
}
