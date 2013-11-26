package textRank

import (
	"bufio"
	r "github.com/ggaaooppeenngg/Gommseg"
	"os"
	"runtime"
	"sort"
)

//获得filePath文档的关键词，这个算法效果不是特别明显，现在主要的分词都是借助机器学习和词典结合的，这里我只用了词典分词
//top是要获得关键词系数最高的top个单词，如果top大于结果大小就直接返回结果
//precision是递归条件，设得大递归时间长，效果好些，但个人感觉0.01 或者0.001都够了，因为后面基本就收敛了
//top key words you want to return
//precision convergence condition
func GetKeyswords(filePath string, top int, precision float64) (rts []string) {
	//初始权重
	var initialWt float64 = 1
	//窗体大小，在同一个窗体内的单词就当作有关系
	//这个可以自己设置
	var window int = 4
	//并发数，根据情况设置
	runtime.GOMAXPROCS(2)
	var s = new(r.Segmenter)
	errInit := s.Init("../Godarts/darts.lib")
	if errInit != nil {
		fmt.Println(errInit)
	}

	offset := 0
	//切分单词
	unifile, _ := os.Open(filePath)
	uniLineReader := bufio.NewReaderSize(unifile, 4000)
	line, bufErr := uniLineReader.ReadString('\n')
	words := make([]string, 0, 100)
	takeWord := func(off int, length int) {
		if len(line[off-offset:off-offset+length]) > 3 {
			words = append(words, string(line[off-offset:off-offset+length]))
		}
	}
	for nil == bufErr {
		s.Mmseg(line[:], offset, takeWord, nil, false)
		offset += len(line)
		line, bufErr = uniLineReader.ReadString('\n')
	}
	s.Mmseg(line, offset, takeWord, nil, false)
	//建立单词对点的映射
	graphMap := make(map[string]Vertex)
	//以点的单词个数建立一个图
	wordsGraph := NewGraph(len(words))
	for _, w := range words {
		_, ok := graphMap[w]
		if !ok {
			//去重
			graphMap[w] = wordsGraph.AddVertex()
		}
	}
	for i := 0; i < len(words)-window; i++ {
		for j := i + 1; j < i+window; j++ {
			//int(1) is the initial weight initialWt
			has := wordsGraph.HasEdge(graphMap[words[i]], graphMap[words[j]])
			if has {
				wordsGraph.AddEdgeWeight(graphMap[words[i]], graphMap[words[j]], initialWt)
				wordsGraph.AddEdgeWeight(graphMap[words[j]], graphMap[words[i]], initialWt)
			} else {
				//感觉无向图更适合单词间的关系,所以正反使用了两次,graph是一个有向图,有点别扭的感觉……
				wordsGraph.AddEdge(graphMap[words[i]], graphMap[words[j]], initialWt)
				wordsGraph.AddEdge(graphMap[words[j]], graphMap[words[i]], initialWt)
			}
		}
	}

	oldScores := make([]float64, wordsGraph.VertexCount(), wordsGraph.VertexCount())
	for i := range oldScores {
		oldScores[i] = 1
	}
	m := float64(10000)
	//run until convergence
	for m > float64(precision) {
		//.85 is the d in TextRank algrithom
		//也就是中文里的阻尼系数
		newScores := Iterate(float64(0.85), oldScores, wordsGraph)

		m = Abs(oldScores, newScores)
		copy(oldScores, newScores)
	}
	pairs := make(IndexScorePairSlice, wordsGraph.VertexCount(), wordsGraph.VertexCount())
	for i, v := range oldScores {
		pairs[i] = IndexScorePair{Vertex(i), v}
	}
	sort.Sort(pairs)
	outMap := make(map[Vertex]string)
	for k, v := range graphMap {
		outMap[v] = k
	}
	if top > len(outMap) {
		rts = make([]string, 0, len(outMap))
	} else {
		rts = make([]string, 0, top)
	}

	for i := 0; i < cap(rts); i++ {
		rts = append(rts, outMap[pairs[i].Index])
	}
	return rts
}
