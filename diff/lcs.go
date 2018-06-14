package diff

import (
	"go/token"
	"strings"

	"github.com/sirupsen/logrus"
)

type simpleColoring struct {
	Color Color
	Data  string
}

func LCS(a, b string, offset int, mode Mode) (coloring Coloring) {
	aList := strings.Split(a, "\n")
	bList := strings.Split(b, "\n")
	c := lcs(aList, bList)
	for _, simpleColor := range printLCS(c, aList, bList, len(aList), len(bList), mode) {
		coloring = append(coloring, ColorChange{Color: simpleColor.Color, Pos: token.Pos(offset), End: token.Pos(offset + len(simpleColor.Data))})
		offset += len(simpleColor.Data) + 1
	}
	return
}

func printLCS(c [][]int, x, y []string, i, j int, mode Mode) (coloring []simpleColoring) {
	if i > 0 && j > 0 && x[i-1] == y[j-1] {
		coloring = printLCS(c, x, y, i-1, j-1, mode)
		coloring = append(coloring, simpleColoring{Color: ColorSame, Data: x[i-1]})
		logrus.Debugln(" " + x[i-1])
	} else {
		if j > 0 && (i == 0 || c[i][j-1] >= c[i-1][j]) {
			coloring = printLCS(c, x, y, i, j-1, mode)
			if mode == ModeNew {
				coloring = append(coloring, simpleColoring{Color: ColorNew, Data: y[j-1]})
				logrus.Debugln("+" + y[j-1])
			}
		} else if i > 0 && (j == 0 || c[i][j-1] < c[i-1][j]) {
			coloring = printLCS(c, x, y, i-1, j, mode)
			if mode == ModeOld {
				coloring = append(coloring, simpleColoring{Color: ColorRemoved, Data: x[i-1]})
				logrus.Debugln("-" + x[i-1])
			}
		}
	}
	return coloring
}

func lcs(x, y []string) [][]int {
	m := len(x)
	n := len(y)
	c := make([][]int, m+1)
	for i := range c {
		c[i] = make([]int, n+1)
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if x[i] == y[j] {
				c[i+1][j+1] = c[i][j] + 1
			} else {
				c[i+1][j+1] = max(c[i+1][j], c[i][j+1])
			}
		}
	}
	return c
}
