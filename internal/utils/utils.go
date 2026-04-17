// Package utils provides utility functions for the Skooma application.
package utils

import "math/rand"

func GetRandomKhajiitPhrase() string {
	messages := []string{
		"🧪 This one is brewing a fresh batch of Skooma...",
		"🦁 Khajiit has wares, if you have coin...",
		"🌙 By Azura! This one crafts magical elixir...",
		"🏝️ May your roads lead you to warm sands...",
		"🧙 This one mixes moon sugar and nightshade...",
		"🏺 Psst! Khajiit knows you come for the good stuff...",
	}
	return messages[rand.Intn(len(messages))]
}
