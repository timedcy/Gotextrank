package textRank

import (
	"fmt"
	"testing"
)

func TestTextRank(t *testing.T) {
	//第一个参数是文件目录，第二个参数是关键词的个数，第三个参数是收敛条件
	rts := GetKeyswords("testArticle.txt", 3, 0.001)
	t.Log(rts)
}
