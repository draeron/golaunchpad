package minimk3

//go:generate go-enum -f=$GOFILE

// Layout x ENUM(
/*
  Pressed
	Released
*/
// )
type EventType int
