package game

type ItemClass int
const (
  Lumber ItemClass = iota
  Stone = iota
)

type Item struct {
  class ItemClass
  quality int
  amount int
}
