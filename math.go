package goutil

// Degree is an integer that rotates 360 degrees.
//
// This method behaves similar to uint8, in how adding anything above 255 will rotate back to 0,
// but instead it will rotate from 360 to 0 (359 max).
type Degree struct {
	deg int16
}

// Deg creates a new Degree integer that rotates 360 degrees.
//
// This method behaves similar to uint8, in how adding anything above 255 will rotate back to 0,
// but instead it will rotate from 360 to 0 (359 max).
func Deg(deg int16) *Degree {
	d := Degree{deg}
	d.clamp()
	return &d
}

// clamp rotates anything above 360 back to 0, and anything below 0 back to 360
func (d *Degree) clamp() {
	for d.deg < 0 {
		d.deg += 360
	}
	for d.deg >= 360 {
		d.deg -= 360
	}
	if d.deg < 0 {
		d.deg = 0
	}
}

// Get rotation value
func (d *Degree) Get() int16 {
	return d.deg
}

// Set new rotation
func (d *Degree) Set(deg int16) {
	d.deg = deg
	d.clamp()
}

// Rotate @deg degrees
func (d *Degree) Rotate(deg int16) {
	d.deg += deg
	d.clamp()
}

// Distance calculates the minimum distance between to angles.
//
// 0 and 15 will return a distance of 15
//
// 0 and 345 will also return a distance of 15, because the rotation from 360 to 0 is smaller
func (d *Degree) Distance(deg *Degree) int16 {
	if d.deg == deg.deg {
		return 0
	}

	d1 := d.deg
	d2 := deg.deg

	n := int16(0)
	if d1 < d2 {
		n = d2 - d1
	} else {
		n = d1 - d2
	}

	d1 = Deg(d1 + (360 / 2)).Get()
	d2 = Deg(d2 + (360 / 2)).Get()

	if d1 < d2 && d2-d1 < n {
		return d2 - d1
	} else if d1 > d2 && d1-d2 < n {
		return d1 - d2
	}

	return n
}
