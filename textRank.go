// textRank project textRank.go
package textRank

import (
	"fmt"
	"math"
)

func Iterate(d float64, oldScores []float64, g *Graph) []float64 {
	ci := make(chan int)
	oldLen := len(oldScores)
	newScores := make([]float64, oldLen, oldLen)
	for i := range oldScores {
		go func(j int) {
			s := float64(0)
			inBound := g.In(Vertex(j))
			for vin, vinVertex := range inBound {
				outBound := g.Out(Vertex(vin))
				outDegree := float64(0)
				for _, out := range outBound {
					outDegree += out.weight
				}
				w, err := g.Weight(vinVertex.index, Vertex(j))
				//IsInf返回表示outDegree是否是正负无限大
				if err != nil {
					fmt.Println(err)
				}
				if outDegree != 0 && !math.IsInf(outDegree, 0) {
					s += (w / outDegree) * oldScores[vinVertex.index]
				}
			}
			newScores[j] = (1 - d) + (d * s)
			ci <- 1
		}(i)
	}
	for j := 0; j < len(oldScores); j++ {
		<-ci
	}
	close(ci)
	return newScores
}

func Abs(s1, s2 []float64) float64 {
	cumulative := float64(0)
	for i, v1 := range s1 {
		cumulative += (v1 - s2[i]) * (v1 - s2[i])
	}
	return cumulative
}

//implement sort interface
//结果对，实现sort的接口
type IndexScorePair struct {
	Index Vertex
	Score float64
}
type IndexScorePairSlice []IndexScorePair

func (s IndexScorePairSlice) Len() int {
	return len(s)
}

func (s IndexScorePairSlice) Swap(i, j int) {
	s[i].Index, s[i].Score, s[j].Index, s[j].Score = s[j].Index, s[j].Score, s[i].Index, s[i].Score
}

//这里反了一下，为了实现 降序输出
//for Desc
func (s IndexScorePairSlice) Less(i, j int) bool {
	if s[i].Score > s[j].Score {
		return true
	} else {
		return false
	}
}
