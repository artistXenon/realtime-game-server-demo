package maze

type Union struct {
	array []int
}

func NewUnion(size int) Union {
	union := Union{ make([]int, size) }
	for i, _ := range union.array {
		union.array[i] = -1
	}
	return union
}

func (u *Union) Find(cell int) int {
	t := 0
	paths := []int{}

	tmpIterator := func() bool {
		t = u.array[cell]
		return t > -1
	}

	for tmpIterator() {
		cell = t
		paths = append(paths, t)
	}

	lp:
	for _, path := range paths {
		if path == cell {
			break lp
		}
		u.array[path] = cell
	}
	return cell
}

func (u *Union) Union(cell1 int, cell2 int) bool {
	i := u.Find(cell1)
	j := u.Find(cell2)
	if i == j {
		return false
	}
	if u.array[i] < u.array[j] {
		k := i
		i = j
		j = k
	}
	u.array[j] += u.array[i]
	u.array[i] = j
	return true
}
