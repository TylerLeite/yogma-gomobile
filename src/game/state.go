package game

type State struct {
  money int
  inventory []Item

  domains []Domain
  buildings []*Building
  populace []Adventurer

  day int

  name string
  title string
}

type Delta struct {
  history []*State
}

func DeltaWith (s1 *State, s2 *State) *Delta {
  return nil
}
