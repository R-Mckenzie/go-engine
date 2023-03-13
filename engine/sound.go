package engine

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var soundBank map[string]beep.Buffer

var playing []beep.StreamSeeker
var looped map[string]beep.Streamer

var isInitialised bool = false

func LoadSound(filepath, name string) {
	if _, ok := soundBank[name]; ok {
		log.Println("Sound already loaded: ", name)
		return
	}
	f, err := os.Open("res/music.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	if !isInitialised {
		soundBank = make(map[string]beep.Buffer)
		playing = make([]beep.StreamSeeker, 0, 10)
		looped = make(map[string]beep.Streamer)
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		isInitialised = true
	}

	buffer := *beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()
	soundBank[name] = buffer
}

func updateSound() {
	toRemove := []int{}
	for i, s := range playing {
		if s.Position() >= s.Len() {
			toRemove = append(toRemove, i)
		}
	}
	for _, i := range toRemove {
		playing[i] = playing[len(playing)-1]
		playing = playing[:len(playing)-1]
	}
}

func LoopSound(name string) {
	sound, ok := soundBank[name]
	if !ok {
		log.Println("sound does not exist: ", name)
		return
	}

	streamer := sound.Streamer(0, sound.Len())
	loop := beep.Loop(-1, streamer)
	looped[name] = loop
	speaker.Play(loop)
}

func PlaySound(name string) {
	sound, ok := soundBank[name]
	if !ok {
		log.Println("sound does not exist: ", name)
		return
	}

	streamer := sound.Streamer(0, sound.Len())
	playing = append(playing, streamer)
	speaker.Play(streamer)
}

func StopLoop(name string) {
	_, ok := looped[name]
	if !ok {
		log.Println("sound is not actively looping: ", name)
		return
	}
	speaker.Lock()
	delete(looped, name)
	speaker.Unlock()
}

func ClearSounds() {
	speaker.Clear()
}

func close() {
	speaker.Close()
}
