package main

type Slider struct {
	max   int
	min   int
	value int
}

func NewSlider(minValue, maxValue, value int) *Slider {
	if value < minValue {
		panic("value is less than min")
	}

	if value > maxValue {
		panic("value is greater than max")
	}

	if minValue > maxValue {
		panic("min is greater than max")
	}

	return &Slider{
		min:   minValue,
		max:   maxValue,
		value: value,
	}
}

func (s *Slider) Up() *Slider {
	if s.value >= s.max {
		return s
	}

	ret := s.clone()
	ret.value++

	return ret
}

func (s *Slider) Down() *Slider {
	if s.value <= s.min {
		return s
	}

	ret := s.clone()
	ret.value--

	return ret
}

func (s *Slider) clone() *Slider {
	return &Slider{
		min:   s.min,
		max:   s.max,
		value: s.value,
	}
}
