package life

import (
	"bufio"
	"errors"
	"math/rand"
	"os"
)

type World struct {
	Height int
	Width  int
	Cells  [][]bool
}

func NewWorld(height, width int) *World {
	cells := make([][]bool, height)
	for i := range cells {
		cells[i] = make([]bool, width)
	}
	return &World{
		Height: height,
		Width:  width,
		Cells:  cells,
	}
}

func (w *World) Seed() {
	for _, row := range w.Cells {
		for i := range row {
			if rand.Intn(10) == 1 {
				row[i] = true
			}
		}
	}
}

func (w *World) IsAlive(x, y int) int {
	if x < 0 || y < 0 || x >= w.Height || y >= w.Width {
		return 0
	}
	if w.Cells[x][y] {
		return 1
	}
	return 0
}

func (w *World) Neighbors(x, y int) int {
	dx := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy := []int{-1, 0, 1, -1, 1, -1, 0, 1}
	cnt := 0
	for i := 0; i < 8; i++ {
		nx := x + dx[i]
		ny := y + dy[i]
		if nx >= 0 && nx < w.Height && ny >= 0 && ny < w.Width && w.Cells[nx][ny] {
			cnt++
		}
	}
	return cnt
}

func (w *World) Next(x, y int) bool {
	n := w.Neighbors(x, y)
	alive := w.Cells[x][y]
	if n <= 3 && n >= 2 && alive {
		return true
	}
	if n == 3 && !alive {
		return true
	}
	return false
}

func NextState(oldWorld, newWorld *World) {
	for i := 0; i < oldWorld.Height; i++ {
		for j := 0; j < oldWorld.Width; j++ {
			newWorld.Cells[i][j] = oldWorld.Next(i, j)
		}
	}
}

func (w *World) SaveState(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()
	for i := 0; i < w.Height; i++ {
		var s string
		for j := 0; j < w.Width; j++ {
			if w.Cells[i][j] {
				s += "1"
			} else {
				s += "0"
			}
		}
		if i != w.Height-1 {
			s += "\n"
		}
		_, err = f.WriteString(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *World) LoadState(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		f.Close()
	}()
	fileScanner := bufio.NewScanner(f)
	width, height := -1, 0
	cells := make([][]bool, 0)
	for fileScanner.Scan() {
		cells = append(cells, make([]bool, 0))
		if width != -1 && width != len(fileScanner.Text()) {
			return errors.New("number of symbols in lines is different")
		}
		width = len(fileScanner.Text())
		for i := 0; i < width; i++ {
			if fileScanner.Text()[i] == '1' {
				cells[height] = append(cells[height], true)
			} else {
				cells[height] = append(cells[height], false)
			}
		}
		height++
	}
	w.Height = height
	w.Width = width
	w.Cells = cells
	return nil
}

func (w *World) String() string {
	brownSquare := "\xF0\x9F\x9F\xAB"
	greenSquare := "\xF0\x9F\x9F\xA9"
	result := ""
	for i := 0; i < w.Height; i++ {
		for j := 0; j < w.Width; j++ {
			if w.Cells[i][j] {
				result += greenSquare
			} else {
				result += brownSquare
			}
		}
		result += "\n"
	}
	return result
}
