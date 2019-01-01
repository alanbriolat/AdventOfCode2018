package util

type GenericSet map[Generic]struct{}

func (s *GenericSet) Add(i Generic) {
	(*s)[i] = struct{}{}
}

func (s *GenericSet) Remove(i Generic) {
	delete(*s, i)
}

func (s *GenericSet) Contains(i Generic) bool {
	_, ok := (*s)[i]
	return ok
}

func (s *GenericSet) Difference(o *GenericSet) {
	for k := range *o {
		s.Remove(k)
	}
}

func (s *GenericSet) Union(o *GenericSet) {
	for k := range *o {
		s.Add(k)
	}
}
