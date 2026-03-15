package main

import (
	"math"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

type soundMsg string

func playProceduralSound(sound string) {
	sr := beep.SampleRate(44100)
	var streamer beep.Streamer

	switch sound {
	case "creak":
		// A low frequency square wave with some decay
		freq := 100.0
		duration := time.Millisecond * 300
		numSamples := int(sr.N(duration))
		i := 0
		streamer = beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
			for j := range samples {
				if i >= numSamples {
					return j, false
				}
				// Square wave
				val := 0.1
				if math.Sin(2*math.Pi*freq*float64(i)/float64(sr)) < 0 {
					val = -0.1
				}
				// Decay
				decay := 1.0 - float64(i)/float64(numSamples)
				samples[j][0] = val * decay
				samples[j][1] = val * decay
				i++
			}
			return len(samples), true
		})
	case "walk":
		// A short low frequency thump
		freq := 50.0
		duration := time.Millisecond * 50
		numSamples := int(sr.N(duration))
		i := 0
		streamer = beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
			for j := range samples {
				if i >= numSamples {
					return j, false
				}
				// Sine wave
				val := math.Sin(2*math.Pi*freq*float64(i)/float64(sr)) * 0.1
				// Decay
				decay := 1.0 - float64(i)/float64(numSamples)
				samples[j][0] = val * decay
				samples[j][1] = val * decay
				i++
			}
			return len(samples), true
		})
	case "splash":
		// White noise with decay
		duration := time.Millisecond * 400
		numSamples := int(sr.N(duration))
		i := 0
		streamer = beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
			for j := range samples {
				if i >= numSamples {
					return j, false
				}
				val := (rand.Float64()*2 - 1) * 0.1
				// Decay
				decay := 1.0 - float64(i)/float64(numSamples)
				samples[j][0] = val * decay
				samples[j][1] = val * decay
				i++
			}
			return len(samples), true
		})
	default:
		return
	}

	speaker.Play(streamer)
}

func playSound(sound string) tea.Cmd {
	return func() tea.Msg {
		playProceduralSound(sound)
		return soundMsg(sound)
	}
}
