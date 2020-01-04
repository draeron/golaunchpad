package device

import (
  "bytes"
  "fmt"
  "github.com/draeron/golaunchpad/pkg/device/event"
  "gitlab.com/gomidi/midi"
  "gitlab.com/gomidi/midi/mid"
  "gitlab.com/gomidi/midi/midimessage/realtime"
  "gitlab.com/gomidi/midi/midireader"
  "gitlab.com/gomidi/midi/midiwriter"
)

type MiniMk3 struct {
  dawIn  mid.In
  midiIn  mid.In
  dawOut  mid.Out
  midiOut mid.Out

  listeners []Listener
}

type Listener func(msg event.Event)

func (m *MiniMk3) Open() error {
  var err error
  err = m.dawIn.Open()
  if err != nil {
    return err
  }

  err = m.dawIn.SetListener(m.receive)
  if err != nil {
    return err
  }

  err = m.dawOut.Open()
  if err != nil {
    return err
  }

  err = m.midiIn.Open()
  if err != nil {
    return err
  }

  err = m.midiIn.SetListener(m.receive)
  if err != nil {
    return err
  }

  err = m.midiOut.Open()
  if err != nil {
    return err
  }

  return nil
}

func (m *MiniMk3) AddListener(listener Listener) {
  m.listeners = append(m.listeners, listener)
}

func (m *MiniMk3) receive(data []byte, deltaMicroseconds int64) {
  rd := midireader.New(bytes.NewReader(data), func(message realtime.Message) {
    log.Debugf(message.String())
  })

  var msg midi.Message
  var err error

  for {
    msg, err = rd.Read()

    // breaking at least with io.EOF
    if err != nil {
      break
    }

    //log.Debugf("msg: %s", msg.String())

    msg, err := event.From(msg)
    if err != nil {
      log.Errorf("%s", err.Error())
    } else if msg.IsValid() {
      for _, l := range m.listeners {
        l(msg)
      }
    }
  }
}

func (m *MiniMk3) SendMidi(msg midi.Message) error {
  return send(m.midiOut, msg)
}

func (m *MiniMk3) SendDaw(msg midi.Message) error {
  return send(m.dawOut, msg)
}

func send(out mid.Out, msg midi.Message) error {
  var bf bytes.Buffer
  wr := midiwriter.New(&bf)
  err := wr.Write(msg)
  if err != nil {
    return err
  }
  return out.Send(bf.Bytes())
}

func (m *MiniMk3) Close() {
  if m.dawIn.IsOpen() {
    m.dawIn.StopListening()
    m.dawIn.Close()
  }

  if m.dawOut.IsOpen() {
    m.dawOut.Close()
  }

  if m.midiIn.IsOpen() {
    m.midiIn.StopListening()
    m.midiIn.Close()
  }

  if m.midiOut.IsOpen() {
    m.midiOut.Close()
  }
}

func (m *MiniMk3) String() string {
  return fmt.Sprintf("Daw In: %s\nDaw Out: %s\nMidi In: %s\nMidi Out: %s\n",
    m.dawIn, m.dawOut, m.midiIn, m.midiOut,
  )
}

func (m *MiniMk3) Print() {
  log.Infof("Daw In: %s", m.dawIn.String())
  log.Infof("Daw Out: %s", m.dawOut.String())
  log.Infof("Midi In: %s", m.midiIn)
  log.Infof("Midi Out: %s", m.midiOut)
}