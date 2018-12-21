// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package audio

import (
	"fmt"
	"io"
	"os"
)

type AudioInfo struct {
	Format     int     // OpenAl Format
	Channels   int     // Number of channels
	SampleRate int     // Sample rate in hz
	BitsSample int     // Number of bits per sample (8 or 16)
	DataSize   int     // Total data size in bytes
	BytesSec   int     // Bytes per second
	TotalTime  float64 // Total time in seconds
}

type AudioFile struct {
	AudioStream
	loops int // replay times
}

type AudioStream interface {
	io.ReadSeeker
	Close() error
	CurrentTime() float64
	Info() AudioInfo
}

// NewAudioFile creates and returns a pointer to a new audio file object and an error
func NewAudioFile(filename string) (*AudioFile, error) {

	// Checks if file exists
	_, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	fin, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return NewAudio(fin)
}

func NewAudio(in io.ReadSeeker) (*AudioFile, error) {
	af := new(AudioFile)

	// Try to open as a wave file
	wav, err := openWave(in)
	if err == nil {
		af.AudioStream = wav
		return af, nil
	}

	in.Seek(0, os.SEEK_SET)
	return nil, fmt.Errorf("Unsuported file type")
}

// Close closes the audiofile
func (af *AudioFile) Close() error {
	return af.AudioStream.Close()
}

// Read reads decoded data from the audio file
func (af *AudioFile) Read(bs []byte) (int, error) {
	nbytes := len(bs)
	n, err := af.AudioStream.Read(bs)
	if err != nil {
		return 0, err
	}
	if n == nbytes {
		return n, nil
	}
	if af.loops == 0 {
		return n, nil
	}
	af.loops--
	// EOF reached. Position file at the beginning
	_, err = af.AudioStream.Seek(0, os.SEEK_SET)
	if err != nil {
		return 0, nil
	}
	// Reads next data into the remaining buffer space
	n2, err := af.AudioStream.Read(bs[n:])
	if err != nil {
		return 0, err
	}
	return n + n2, err
}

// Looping returns the current looping state of this audio file
func (af *AudioFile) Looping() bool {

	return af.loops > 0
}

// SetLooping sets the looping state of this audio file
func (af *AudioFile) SetLooping(loops int) {

	af.loops = loops
}
