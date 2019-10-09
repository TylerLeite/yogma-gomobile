package game

type AdventurerClass int
const (
  Fighter AdventurerClass = iota
  Wizard = iota
  Sorcerer = iota
  Archer = iota
  Barbarian = iota
  Thief = iota
  Druid = iota
  Cleric = iota
  Bard = iota
  Psion = iota
  Monk = iota
  Artificer = iota
  Necromancer = iota
)

type AdventurerState int
const (
  Idle AdventurerState = iota
  Roaming = iota
  Shopping = iota
  Fighting = iota
  Questing = iota
  Leaving = iota
)

type Adventurer struct {
  class AdventurerClass
  prestigious bool
  level int

  state AdventurerState
  currentObjective *Building // Where is this adventurer headed

  thoughts []thought

  money int
  inventory []Item

  // Location
  currentDomain *Domain
  x,y int
}

type thought struct {
  subject string // usually a building
  object string // e.g. price
  positivity int // -100 to 100
}
