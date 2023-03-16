package engine

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type gameSound struct {
	ctrl *beep.Ctrl
	vol  *effects.Volume
}

var soundBank map[string]beep.Buffer

var looped map[string]*gameSound

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

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	if !isInitialised {
		soundBank = make(map[string]beep.Buffer)
		looped = make(map[string]*gameSound)
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
	vol := &effects.Volume{Streamer: streamer, Base: 2, Volume: volume, Silent: false}
	speaker.Play(vol)
}

func LoopSound(name string, volume float64) {
	fmt.Println("Loop Sound")
	sound, ok := soundBank[name]
	if !ok {
		log.Println("sound does not exist: ", name)
		return
	}

	existing, ok := looped[name]
	if !ok {
		fmt.Println("!exists")
		streamer := sound.Streamer(0, sound.Len())
		loop := beep.Loop(-1, streamer)
		ctrl := &beep.Ctrl{Streamer: loop, Paused: false}
		vol := &effects.Volume{
			Streamer: ctrl,
			Base:     10,
			Volume:   volume,
			Silent:   false,
		}
		looped[name] = &gameSound{ctrl: ctrl, vol: vol}
		speaker.Play(looped[name].vol)
	} else {
		fmt.Println("exists")
		speaker.Lock()
		existing.ctrl.Paused = false
		speaker.Unlock()
	}
}

func StopLoop(name string) {
	fmt.Println("Stop Loop")
	loop, ok := looped[name]
	if !ok {
		log.Println("sound is not actively looping: ", name)
		return
	}
	speaker.Lock()
	loop.ctrl.Paused = true
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
	loop.ctrl.Paused = true
	speaker.Unlock()
}

func ClearSounds() {
	speaker.Clear()
}

func close() {
	ClearSounds()
	speaker.Close()
}
