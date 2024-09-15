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
	onChange           []func(newSong *subsonic.Child, isPaused bool)
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

func (q *Queue) Insert(i int, s *subsonic.Child) {
	if len(q.songs) == 0 {
		q.Append(s)
		return
	}
	q.songs = append(q.songs[:i], append([]*subsonic.Child{s}, q.songs[i:]...)...)
}

func (q *Queue) Clear() {
	q.songs = make([]*subsonic.Child, 0)
	if q.isPaused {
		q.TogglePause()
	}
	q.Stop()
	q.triggerChange()
}

func (q *Queue) PlaySong(s *subsonic.Child) error {
	if q.isPaused {
		q.TogglePause()
	}

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

	q.triggerChange()

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
		for _, f := range q.onChange {
			f(nil, false)
		}
		return nil
	}

	return q.Play()
}

func (q *Queue) Stop() {
	if q.isPaused {
		q.TogglePause()
	}
	speaker.Clear()
}

func (q *Queue) SetOnChangeCallback(f func(newSong *subsonic.Child, isPlaying bool)) {
	q.onChange = append(q.onChange, f)
}

func (q *Queue) TogglePause() {
	if q.isPaused {
		speaker.Unlock()
	} else {
		speaker.Lock()
	}

	q.isPaused = !q.isPaused

	q.triggerChange()
}

func (q *Queue) SkipTo(s *subsonic.Child) {
	i := -1
	for n, s2 := range q.GetSongs() {
		if s.ID == s2.ID {
			i = n
			break
		}
	}

	if i == -1 {
		return
	}

	q.songs = q.songs[i:]
	q.Play()
}

func (q *Queue) RemoveSong(i int) error {
	if i >= len(q.songs) {
		return fmt.Errorf("index out of bounds")
	}

	q.songs = append(q.songs[:i], q.songs[i+1:]...)
	if i == 0 {
		// We removed the first song: this stops it and prepares for the next
		q.Stop()
		if !q.isPaused {
			q.Play()
		}
	}
	q.triggerChange()

	return nil
}

func (q *Queue) Switch(a, b int) error {
	if a >= len(q.songs) {
		return fmt.Errorf("%d is out of bounds", a)
	}

	if b >= len(q.songs) {
		return fmt.Errorf("%d is out of bounds", b)
	}

	tmp := q.songs[a]
	q.songs[a] = q.songs[b]
	q.songs[b] = tmp

	if (a == 0 || b == 0) && !q.isPaused {
		// If we're switching the first song, and it's currently playing, start Play() again
		q.Play()
	}

	q.triggerChange()

	return nil
}

func (q *Queue) triggerChange() {
	if len(q.songs) > 0 {
		for _, f := range q.onChange {
			f(q.songs[0], q.isPaused)
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
		sr := p.oldSampleRate
		p.oldSampleRate = format.SampleRate
		return beep.Resample(4, format.SampleRate, sr, s), nil
	}
}
