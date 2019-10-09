package game

type DomainClass int
const (
  Plains DomainClass = iota
  Hills = iota
  Mountains = iota
  Desert = iota
  Swamp = iota
  Lake = iota
  Coast = iota
  Sea = iota
)

// In tiles
const DomainWidth int = 32
const DomainHeight int = 32
const DomainBytes int = 128 // 32 * 32 / 8

type Domain struct {
  terrain DomainClass
  neighbors []*Domain

  buildings []Building

  x,y int
}

func DomainCoordsToAbsolute (x int, y int) (u int, v int) {
  // TODO
  return 0, 0
}
