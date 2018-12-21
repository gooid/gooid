// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package audio

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/gooid/audio/al"
)

// WaveSpecs describes the characterists of the audio encoded in a wave file.
type AudioInfo_ struct {
	Format     int     // OpenAl Format
	Type       int     // Type field from wave header
	Channels   int     // Number of channels
	SampleRate int     // Sample rate in hz
	BitsSample int     // Number of bits per sample (8 or 16)
	DataSize   int     // Total data size in bytes
	BytesSec   int     // Bytes per second
	TotalTime  float64 // Total time in seconds
}

const (
	waveHeaderSize = 44
	fileMark       = "RIFF"
	fileHead       = "WAVE"
)

// WaveCheck checks if the specified filepath corresponds to a an audio wave file.
// If the file is a valid wave file, return a pointer to WaveSpec structure
// with information about the encoded audio data.
func WaveCheckFile(filepath string) (*AudioInfo, error) {

	// Open file
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return WaveCheck(f)
}

func WaveCheck(f io.Reader) (*AudioInfo, error) {
	// Reads header
	header, err := readWavHead(f)
	if err != nil {
		return nil, err
	}

	// Decodes header fields
	var ws AudioInfo
	//ws.Type = int(header.Format)
	ws.Channels = int(header.Channels)
	ws.SampleRate = int(header.SamplesPerSec)
	ws.BitsSample = int(header.BitsPerSample)
	ws.DataSize = int(header.DataSize)

	// Sets OpenAL format field if possible
	if ws.Channels == 1 {
		if ws.BitsSample == 8 {
			ws.Format = al.FormatMono8
		} else if ws.BitsSample == 16 {
			ws.Format = al.FormatMono16
		}
	} else if ws.Channels == 2 {
		if ws.BitsSample == 8 {
			ws.Format = al.FormatStereo8
		} else if ws.BitsSample == 16 {
			ws.Format = al.FormatStereo16
		}
	}

	// Calculates bytes/sec and total time
	bytesChannel := ws.BitsSample / 8

	ws.BytesSec = ws.SampleRate * ws.Channels * bytesChannel
	ws.TotalTime = float64(ws.DataSize) / float64(ws.BytesSec)
	return &ws, nil
}

type wavHeader struct {
	Riff          [4]byte //'RIFF'
	RiffSize      uint32  //44+DataSize
	Wave          [4]byte //'WAVE'
	Fmt           [4]byte //'fmt '
	FmtSize       uint32  // 16
	Format        uint16  // 1
	Channels      uint16
	SamplesPerSec uint32
	BytesPerSec   uint32
	BlockAlign    uint16
	BitsPerSample uint16
	Data          [4]byte //'data'
	DataSize      uint32
}

func readWavHead(r io.Reader) (*wavHeader, error) {
	var head wavHeader
	err := binary.Read(r, binary.LittleEndian, &head)
	if err != nil {
		return nil, err
	}
	// Checks file marks
	if string(head.Riff[:]) != fileMark {
		return nil, fmt.Errorf("'RIFF' mark not found")
	}
	if string(head.Wave[:]) != fileHead {
		return nil, fmt.Errorf("'WAVE' mark not found")
	}
	//"fmt " == string(head.Fmt[:])
	for "data" != string(head.Data[:]) {
		_, err := r.Read(head.Data[:])
		if err != nil {
			return nil, err
		}
		if "data" == string(head.Data[:]) {
			binary.Read(r, binary.LittleEndian, &head.DataSize)
			return &head, nil
		}
	}
	return &head, nil
}

type WavFile struct {
	reader io.ReadSeeker
	offset int64 // data offset
	info   AudioInfo
}

func openWave(in io.ReadSeeker) (*WavFile, error) {
	info, err := WaveCheck(in)
	if err != nil {
		return nil, err
	}
	if info.Format == -1 {
		return nil, fmt.Errorf("Unsupported OpenAL format")
	}

	pos, _ := in.Seek(0, os.SEEK_CUR)
	return &WavFile{
		reader: in,
		offset: pos,
		info:   *info,
	}, nil
}

func (w *WavFile) Read(p []byte) (int, error) {
	return w.reader.Read(p)
}

func (w *WavFile) Seek(offset int64, whence int) (int64, error) {
	if whence == os.SEEK_SET {
		return w.reader.Seek(offset+w.offset, whence)
	} else {
		return w.reader.Seek(offset, whence)
	}
}

func (w *WavFile) CurrentTime() float64 {
	pos, err := w.reader.Seek(0, os.SEEK_CUR)
	if err != nil {
		return 0
	}
	return float64(pos-w.offset) / float64(w.info.BytesSec)
}

func (w *WavFile) Info() AudioInfo {
	return w.info
}

func (w *WavFile) Close() error {
	return nil
}

func WavWriteHeader(w io.Writer, sampleRate, format int) error {
	head := wavHeader{
		Riff:    [4]byte{'R', 'I', 'F', 'F'}, //'RIFF'
		Wave:    [4]byte{'W', 'A', 'V', 'E'}, //'WAVE'
		Fmt:     [4]byte{'f', 'm', 't', ' '}, //'fmt '
		FmtSize: 16,                          // 16
		Format:  1,                           // 1
		Data:    [4]byte{'d', 'a', 't', 'a'}, //'data'
	}

	switch format {
	case al.FormatMono8:
		head.Channels = 1
		head.BitsPerSample = 8
	case al.FormatMono16:
		head.Channels = 1
		head.BitsPerSample = 16
	case al.FormatStereo8:
		head.Channels = 2
		head.BitsPerSample = 8
	case al.FormatStereo16:
		head.Channels = 2
		head.BitsPerSample = 16
	}
	head.SamplesPerSec = uint32(sampleRate)
	head.BlockAlign = head.Channels * head.BitsPerSample / 8
	head.BytesPerSec = head.SamplesPerSec * uint32(head.Channels) * uint32(head.BitsPerSample) / 8

	return binary.Write(w, binary.LittleEndian, &head)
}

func WavClose(w io.WriteSeeker) error {
	offset64, err := w.Seek(0, os.SEEK_CUR)
	offset := uint32(offset64)
	if err == nil {
		w.Seek(4, os.SEEK_SET)
		binary.Write(w, binary.LittleEndian, offset)
		w.Seek(40, os.SEEK_SET)
		binary.Write(w, binary.LittleEndian, offset-44)
	}
	return err
}
