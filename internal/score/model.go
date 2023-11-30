package score

type Score struct {
	user    uint
	total   int
	version uint
}

type View struct {
	User  uint
	Total int
}

func NewView(score *Score) View {
	return View{
		User:  score.user,
		Total: score.total,
	}
}

func NewScore(user uint, total int) *Score {
	return &Score{
		user:    user,
		total:   total,
		version: 0,
	}
}

func (s *Score) AddScore(value int) {
	s.total += value
}

func (s *Score) Total(total int) {
	s.total = total
}
