package maze

import (
	"math/rand"
	"math"
	// "fmt"
	"encoding/json"
)


type Maze struct {
	width int
	height int
	weight [9]float64
	Data *[]int

	wallqs [9]*[]wallPosition
	union *Union
}

type wallPosition struct {
	x int
	y int
	v bool
}


func NewMaze() *Maze {
	maze := Maze{width: 27, height: 18}
	maze.SetWeight([9]int{
		12, 10, 10, 
		10, 10, 8, 
		10, 8, 2})
	return &maze
}

func (m *Maze) SetWidth(w int) *Maze {
	m.width = w
	return m
}

func (m *Maze) SetHeight(h int) *Maze {
	m.height = h
	return m
}

func (m *Maze) SetWeight(w [9]int) *Maze {
	m.weight = [9]float64{}
	for i := 0; i < 9; i++ {
		m.weight[i] = math.Pow(2, float64(w[i]) / 2)
	}
	return m
}

func (m *Maze) SetData(d *[]int) *Maze {
	m.Data = d
	return m
}


func (m *Maze) Init() *Maze {
	data := make([]int, m.width * m.height * 2)
	m.Data = &data
	return m
}


// private

func (m *Maze) setWall(pos wallPosition, wall int) {
	if pos.x >= 0 && pos.x < m.width && pos.y >= 0 && pos.y < m.height {
		vertical := 1 
		if pos.v {
			vertical = 0
		}
		
		(*m.Data)[(pos.y * m.width + pos.x) * 2 + vertical] = wall
	}
}

func (m *Maze) getWall(pos wallPosition) int {
	if pos.x >= 0 && pos.x < m.width && pos.y >= 0 && pos.y < m.height {
		vertical := 1 
		if pos.v {
			vertical = 0
		}
		return (*m.Data)[(pos.y * m.width + pos.x) * 2 + vertical]
	}
	return -1
}

func (m *Maze) getNeighbours(pos wallPosition) [6]wallPosition {
	if pos.v {
		return [6]wallPosition {
			wallPosition{ x: pos.x, y: pos.y - 1, v: false },
			wallPosition{ x: pos.x, y: pos.y - 1, v: true },
			wallPosition{ x: pos.x + 1, y: pos.y - 1, v: false },
			wallPosition{ x: pos.x, y: pos.y, v: false },
			wallPosition{ x: pos.x, y: pos.y + 1, v: true },
			wallPosition{ x: pos.x + 1, y: pos.y, v: false },
			}
	} else {
		return [6]wallPosition {
			wallPosition{ x: pos.x - 1, y: pos.y, v: true },
			wallPosition{ x: pos.x - 1, y: pos.y, v: false },
			wallPosition{ x: pos.x - 1, y: pos.y + 1, v: true },
			wallPosition{ x: pos.x, y: pos.y, v: true },
			wallPosition{ x: pos.x + 1, y: pos.y, v: false },
			wallPosition{ x: pos.x, y: pos.y + 1, v: true },
		}
	}
}

func (m *Maze) getWallEndType(pos1 wallPosition, pos2 wallPosition, pos3 wallPosition) int {
	c := [3]bool { m.getWall(pos1) == 0, m.getWall(pos2) == 0, m.getWall(pos3) == 0 }
	if c[0] == c[2] {
		if c[0] {
			return 2
		} else {
			return 0
		}
	} else {
		if c[1] {
			return 2
		} else {
			return 1
		}
	}
}

func (m *Maze) getWallType(pos wallPosition) int {
	wn := m.getNeighbours(pos)
	t1 := m.getWallEndType(wn[0], wn[1], wn[2])
	t2 := m.getWallEndType(wn[3], wn[4], wn[5])
	return t1 * 3 + t2
}

func (m *Maze) Generate(c chan int) {
	newQueue := [9]*[]wallPosition {
		&[]wallPosition{}, &[]wallPosition{}, &[]wallPosition{}, 
		&[]wallPosition{}, &[]wallPosition{}, &[]wallPosition{}, 
		&[]wallPosition{}, &[]wallPosition{}, &[]wallPosition{},
	}

	m.wallqs = newQueue

	vwall := 0
	hwall := 0
	

	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ { 
			base_index := (y * m.width + x) * 2
			vwall = -1
			hwall = -1
			if x < m.width - 1 {
				vwall = 1
				randominsert(m.wallqs[0], wallPosition{ x: x, y: y, v: true })
			}
			
			if y < m.height - 1 {
				hwall = 1
				randominsert(m.wallqs[0], wallPosition{ x: x, y: y, v: false })
			}

			(*m.Data)[base_index] = vwall
			(*m.Data)[base_index + 1] = hwall
		}
	}

	union := NewUnion(len(*m.Data) / 2)
	m.union = &union

	w := float64(0)
	chunks := make([]float64, 9)

	perf := 0;

	mazemain:
	for true {
		// if perf % 5 == 0 {
		// 	fmt.Printf("%d %d %d %d %d %d %d %d %d\n", 
		// 	len(*m.wallqs[0]), len(*m.wallqs[1]), len(*m.wallqs[2]), 
		// 	len(*m.wallqs[3]), len(*m.wallqs[4]), len(*m.wallqs[5]), 
		// 	len(*m.wallqs[6]), len(*m.wallqs[7]), len(*m.wallqs[8]))
		// }
		perf++
		w = 0
		for i := 0; i < 9; i++ {
			chunks[i] = float64(len(*m.wallqs[i])) * m.weight[i]
			w += chunks[i]
		}

		if w == 0 {
			break mazemain
		}

		w = rand.Float64() * w
		
		i := 0
		qselect:
		for i = 0; i < 8; i++ {
			chunk := chunks[i]
			if w < chunk {
				break qselect
			}
			w -= chunk
		}

		pos := (*m.wallqs[i])[len(*m.wallqs[i]) - 1]
		*m.wallqs[i] = (*m.wallqs[i])[:len(*m.wallqs[i]) - 1]
		wall := m.getWall(pos)
		wallType := m.getWallType(pos)

		if wall == 1 && wallType == i {
			c1 := pos.y * m.width + pos.x
			c2 := c1
			if pos.v {
				c2 += 1
			} else {
				c2 += m.width
			}
			// fmt.Printf("what have i done\n")

			if !m.union.Union(c1, c2) {
				continue mazemain
			}
			neighbours := m.getNeighbours(pos)
			neighbourLp:
			for j := 0; j < 6; j++ {
				m.setWall(pos, 1)
				nType := m.getWallType(neighbours[j])
				m.setWall(pos, 0)
				neighbour := m.getWall(neighbours[j])
				if neighbour == -1 {
					continue neighbourLp
				}
				typ := m.getWallType(neighbours[j])
				if typ <= nType {
					continue neighbourLp
				}
				randominsert(m.wallqs[typ], neighbours[j])
			}	
		}
	}
	c <- 0
}

func (m *Maze) Serialize() string {
	data := make(map[string]interface{})
	data["size"] = m.width
	data["data"] = m.Data
	data["weight"] = m.weight

	doc, _ := json.Marshal(data) // 맵을 JSON 문서로 변환

	return string(doc)
}

// func Deserialize(str string) *Maze {
// 	data := map[string]interface{} 
// 	json.Unmarshal([]byte(str), &data)
// 	maze := NewMaze()
// 	width := (int) data["size"]
// 	mazeData := ([]int) data["data"]
// 	maze.SetWidth(width).SetHeight(int(len(mazeData) / width)).SetData(&mazeData)
// 	fmt.Println(data["name"], data["age"])
// }


func randominsert(arr *[]wallPosition, item wallPosition) {
	if len(*arr) == 0 {
		a := append(*arr, item)
		*arr = a
		return
	}
	i := rand.Intn(len(*arr))

	a := append((*arr)[:i + 1], (*arr)[i:]...)
	a[i] = item
	*arr = a
}