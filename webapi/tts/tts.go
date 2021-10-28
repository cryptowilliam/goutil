package tts

import (
	"github.com/nubunto/tts"
)

// This API needs Google translate access
func Speak(text string, language string) (mp3 []byte, err error) {
	sound, err := tts.Speak(tts.Config{
		Speak:    text,
		Language: "pt-BR",
	})
	if err != nil {
		return nil, err
	}

	return sound.Bytes(), nil
}
