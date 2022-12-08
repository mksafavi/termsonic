package music

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/delucks/go-subsonic"
	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
)

type Queue struct {
	songs    []*subsonic.Child
	isPaused bool

	sub                *subsonic.Client
	speakerInitialized bool
	oldSampleRate      beep.SampleRate
	onChange           func(newSong *subsonic.Child, isPaused bool)
}

func NewQueue(client *subsonic.Client) *Queue {
	return &Queue{
		sub:                client,
		speakerInitialized: false,
	}
}

func (q *Queue) SetClient(client *subsonic.Client) {
	q.Clear()
	q.sub = client
}

func (p *Queue) GetSongs() []*subsonic.Child {
	return p.songs
}

func (q *Queue) Append(s *subsonic.Child) {
	q.songs = append(q.songs, s)
}

func (q *Queue) Clear() {
	q.songs = make([]*subsonic.Child, 0)
	speaker.Clear()
}

func (q *Queue) PlaySong(s *subsonic.Child) error {
	rc, err := Download2(q.sub, s.ID)
	if err != nil {
		return err
	}

	var ssc beep.StreamSeekCloser
	var format beep.Format

	switch filepath.Ext(s.Path) {
	case ".mp3":
		ssc, format, err = mp3.Decode(rc)
	case ".ogg":
		fallthrough
	case ".oga":
		ssc, format, err = vorbis.Decode(rc)
	case ".flac":
		ssc, format, err = flac.Decode(rc)
	default:
		return fmt.Errorf("unknown file type: %s", filepath.Ext(s.Path))
	}

	if err != nil {
		return err
	}

	streamer, err := q.setupSpeaker(ssc, format)
	if err != nil {
		return err
	}
	speaker.Clear()
	speaker.Play(beep.Seq(streamer, beep.Callback(func() { go q.Next() })))

	if q.onChange != nil {
		q.onChange(s, false)
	}

	return nil
}

func (q *Queue) Play() error {
	if len(q.songs) == 0 {
		return fmt.Errorf("the queue is empty")
	}

	s := q.songs[0]
	q.PlaySong(s)

	return nil
}

func (q *Queue) Next() error {
	q.Stop()

	if len(q.songs) == 0 {
		return nil
	}

	q.songs = q.songs[1:]

	if len(q.songs) == 0 {
		if q.onChange != nil {
			q.onChange(nil, false)
		}
		return nil
	}

	return q.Play()
}

func (q *Queue) Stop() {
	speaker.Clear()
}

func (q *Queue) SetOnChangeCallback(f func(newSong *subsonic.Child, isPlaying bool)) {
	q.onChange = f
}

func (q *Queue) TogglePause() {
	if q.isPaused {
		speaker.Unlock()
	} else {
		speaker.Lock()
	}

	q.isPaused = !q.isPaused

	if q.onChange != nil {
		if len(q.songs) > 0 {
			q.onChange(q.songs[0], q.isPaused)
		}
	}
}

func (p *Queue) setupSpeaker(s beep.Streamer, format beep.Format) (beep.Streamer, error) {
	if !p.speakerInitialized {
		err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err != nil {
			return nil, fmt.Errorf("speaker.Init: %v", err)
		}
		p.speakerInitialized = true
		p.oldSampleRate = format.SampleRate

		return s, nil
	} else {
		return beep.Resample(4, format.SampleRate, p.oldSampleRate, s), nil
	}
}
