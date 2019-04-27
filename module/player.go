package module

type Player struct {
	ID         string
	Position   Position
	Velocity   Velocity
	Control    Control
	Class      string
	Attributes Attributes
}

type Position struct {
	X float64
	Y float64
}

type Velocity struct {
	X float64
	Y float64
}

type Control struct {
	Right bool
	Left  bool
	Up    bool
	Down  bool
}

type Attributes struct {
	Attack    int
	Intellect int
	Defense   int
}

func (p *Player) UpdatePosition() {

	if p.Control.Up {
		p.Position.Y -= p.Velocity.Y
	}
	if p.Control.Down {
		p.Position.Y += p.Velocity.Y
	}
	if p.Control.Left {
		p.Position.X -= p.Velocity.X
	}
	if p.Control.Right {
		p.Position.X += p.Velocity.X
	}

}
