package utils

import (
	"math/rand"
	"time"
)

func GenerateToken(minutes int) (int, int) {
	src := rand.NewSource(time.Now().Unix())
	rnd := rand.New(src)
	token := rnd.Intn(999999-100000+1) + 100000
	expiresAt := time.Now().Add(time.Minute).Second() * 1000
	return token, expiresAt
}

func Map[T, V comparable](data []T, fn func(v T) V) []V {
	result := []V{}
	for _, dt := range data {
		result = append(result, fn(dt))
	}
	return result
}

func ForEach[T any](data []T, fn func(v T)) {
	for _, dt := range data {
		fn(dt)
	}
}

func Filter[T comparable](data []T, fn func(v T) bool) []T {
	result := []T{}
	for _, dt := range data {
		if fn(dt) {
			result = append(result, dt)
		}
	}
	return result
}

func LoopOverMap[T comparable](data map[T]T, fn func(k, V T)) {
	for k, v := range data {
		fn(k, v)
	}
}
