package stats

type Repository interface {
	Inc(Key)
	Top() (Top, bool)
}
