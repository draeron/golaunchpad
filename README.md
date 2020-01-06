# golang library for Novation's Launchpad

## Features

- [x] coloring buttons
- [x] button identification
- [x] Sleep / Wake
- [x] auto detect
- [x] programmer mode
- [x] layout (drum, keys, users...) selection
- [ ] support for multiple launchpad
- [ ] Add support for other launchpad devices (x, pro etc)
- [ ] support for daw fader banks / setup
- [ ] LED Feedbacks 
- [ ] Session color
- [ ] Pulse / Flash color
- [ ] Midi clock

## Bugs

- [ ] device inquiry
- [ ] SysEx message readback


## Examples

- [wave](examples/wave/wave.go): colored wave upon pressing a pad button
- [scan](examples/scan/scan.go): light up every button one after the other
- [text](examples/scan/text.go): display pressed coordinate as text

## Ref

- [Launchpad Mini MK3 Programmer’s reference manual](ref/Launchpad_Mini_Programmers_Reference_Manual.pdf)
- [gomidi](https://gitlab.com/gomidi/midi)
