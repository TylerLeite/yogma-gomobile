package game

type BuildingClass int
const (
  cottage BuildingClass = iota
)

type Building struct {
  level int
  _type int

  parent *Domain

  inventory []Item
  money int

  x0,y0, x1,y1 int
}

func (b *Building) IsShop () bool {
  // TODO
  return true;
}

func (b *Building) PayTaxes (n int) int {
  if (b.money >= n) {
    b.money -= n;
  } else {
    n = b.money;
    b.money = 0;
  }

  return n;
}
