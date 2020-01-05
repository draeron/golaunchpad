package event

//go:generate go-enum -f=$GOFILE --noprefix

// Type x ENUM(
/*
	Pressed
	Released
*/
// )
type Type int
