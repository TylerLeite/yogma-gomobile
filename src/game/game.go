package game

import (
  "math/rand"
)

// How many frames per day
const GameSpeed int = 30
var DayCounter int = 0

func calculateTaxes (n int) int {
  if n > 4 {
    return n*5
  } else {
    return n*n
  }
}

func RunOneDay (frame int, s *State) {
  if frame % GameSpeed != 0 {
    return
  }

  // Run adventurer ai (including quests)
  for range s.populace {
    var d *Domain = &s.domains[rand.Intn(len(s.domains))]
    var building *Building = &d.buildings[rand.Intn(len(s.domains))]

    // TODO: purchase with purpose
    if rand.Intn(100) > 33 {
      building.money += rand.Intn(10)
    }
  }

  // Advance the world 1 domain at a time
  for _, domain := range s.domains {
    // Check for events

    // Progress buildings
    const taxAmt int = 1
    var taxes int = 0
    for _, building := range domain.buildings {
      // Collect taxes
      taxes += building.PayTaxes(calculateTaxes(len(s.domains)))
    }
  }

  // Run monster ai


  // Increment the day
  DayCounter++
}
