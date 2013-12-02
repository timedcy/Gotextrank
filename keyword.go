package textRank

import (
	"bufio"
	"fmt"
	r "github.com/ggaaooppeenngg/Gommseg"
	"os"
	"sort"
)

var segmenter *r.Segmenter
var window int = 4
var initialWt float64 = 1.0

func init() {
	segmenter = new(r.Segmenter)
	errInit := segmenter.Init("../Godarts/darts.lib")
	if errInit != nil {
		fmt.Println(errInit)
	}
}

func builtGraph(words []string) (map[string]Vertex, *Graph) {
	graphMap := make(map[string]Vertex)
	wordsGraph := NewGraph(len(words))
	for _, w := range words {
		_, ok := graphMap[w]
		if !ok {
			graphMap[w] = wordsGraph.AddVertex()
		}
	}
	//fmt.Println(graphMap)
	for i := 0; i < len(words)-window; i++ {
		for j := i + 1; j < i+window; j++ {
			from := graphMap[words[i]]
			to := graphMap[words[j]]
			//fmt.Println(from, to)
			has := wordsGraph.HasEdge(from, to)
			if has {
				wordsGraph.AddEdgeWeight(from, to, initialWt)
				wordsGraph.AddEdgeWeight(to, from, initialWt)
			} else {
				//感觉无向图更适合单词间的关系,所以正反使用了两次,graph是一个有向图,有点别扭的感觉……
				wordsGraph.AddEdge(from, to, initialWt)
				wordsGraph.AddEdge(to, from, initialWt)
			}
		}
	}
	//fmt.Println("go")
	//fmt.Println(wordsGraph)
	//fmt.Println(graphMap)
	return graphMap, wordsGraph
}
func iterate(graphMap map[string]Vertex, wordsGraph *Graph, precision float64) (rts []string) {
	m := float64(10000)
	count := wordsGraph.VertexCount()
	oldScores := make([]float64, count, count)
	for i := range oldScores {
		oldScores[i] = 1
	}
	for m > precision {
		newScores := Iterate(float64(0.85), oldScores, wordsGraph)
		//fmt.Println(oldScores)
		m = Abs(oldScores, newScores)
		copy(oldScores, newScores)
	}
	pairs := make(IndexScorePairSlice, count, count)
	for i, v := range oldScores {
		pairs[i] = IndexScorePair{Vertex(i), v}
	}
	sort.Sort(pairs)
	outMap := make(map[Vertex]string)
	for k, v := range graphMap {
		outMap[v] = k
	}
	rts = make([]string, 0, len(outMap))
	for i := 0; i < len(outMap); i++ {
		rts = append(rts, outMap[pairs[i].Index])
	}
	return rts
}
func GetKeyWords(input string, top int, precision float64) (rts []string) {
	words := make([]string, 0, 100)
	takeWord := func(offset, length int) {
		if len(input[offset:offset+length]) > 3 {
			words = append(words, input[offset:offset+length])
			//fmt.Println("go", offset, length, input[offset:offset+length])
		}
	}
	segmenter.Mmseg(input, 0, takeWord, false)
	graphMap, wordsGraph := builtGraph(words)
	//fmt.Println(words)
	results := iterate(graphMap, wordsGraph, precision)
	if len(results) < top {
		top = len(results) - 1
	}
	rts = results[0:top]
	return rts
}

func GetKeyWordsFile(filePath string, top int, precision float64) (rts []string) {
	offset := 0
	unifile, _ := os.Open(filePath)
	uniLineReader := bufio.NewReaderSize(unifile, 4000)
	line, bufErr := uniLineReader.ReadString('\n')
	words := make([]string, 0, 100)
	takeWord := func(off int, length int) {
		if len(line[off-offset:off-offset+length]) > 3 {
			words = append(words, string(line[off:off+length]))
			//fmt.Println(off, length, line[off-offset:off-offset+length])
		}
	}
	for nil == bufErr {
		segmenter.Mmseg(line[:], offset, takeWord, false)
		offset += len(line)
		line, bufErr = uniLineReader.ReadString('\n')
	}
	segmenter.Mmseg(line, offset, takeWord, false)
	graphMap, wordsGraph := builtGraph(words)
	//fmt.Println(words)
	results := iterate(graphMap, wordsGraph, precision)
	if len(results) < top {
		top = len(results) - 1
	}
	rts = results[0:top]
	return rts
}
