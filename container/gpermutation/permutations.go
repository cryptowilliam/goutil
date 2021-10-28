package gpermutation

// Permutations(排列组合)
// 不同种类间有顺序、多分类必选且只选一
// 用于罗列在允许范围内所有可能的参数集，比如量化
// AddKind("a")
// AddKind("b", "c")
// AddKind("d", "e")
// ListAll() => [["a","b","d"], ["a","b","e"], ["a","c","d"], ["a","c","e"]]

// reference: https://www.zhangshengrong.com/p/7B1LBqRNwp/

import (
	"github.com/cryptowilliam/goutil/sys/gcounter"
	"sync/atomic"
)

type Permutations struct {
	list     [][]interface{}
	currLoop atomic.Value
}

type Result struct {
	Items []interface{}
}

func New() *Permutations {
	return &Permutations{}
}

func (c *Permutations) AddKind(items ...interface{}) {
	var itemsArray []interface{}
	itemsArray = append(itemsArray, items...)
	c.list = append(c.list, itemsArray)
}

// 获取所有组合的总数
func (c *Permutations) GetCombinationsCount() int {
	rstCount := 1
	for _, v := range c.list {
		rstCount *= len(v)
	}
	return rstCount
}

func (c *Permutations) GetCurrentLoop() int {
	return -1
}

func (c *Permutations) ListAll() []Result {
	// 计算结果条目总数
	rstCount := 1
	for _, v := range c.list {
		rstCount *= len(v)
	}

	// 计算每条结果的元素长度
	itemCount := len(c.list)

	// 缓存，每行对应一个Kind
	// get error: non-constant array bound itemCount
	//cache := [rstCount][itemCount]interface{}{}
	cache := make([][]interface{}, rstCount)
	for i := range cache {
		cache[i] = make([]interface{}, itemCount)
	}

	// 生成全部cache
	for i := 0; i < itemCount; i++ {

		// 计算步进
		step := 1
		for k := i + 1; k < len(c.list); k++ {
			step *= len(c.list[k])
		}

		// 根据步进创建步进累加器
		sa, err := gcounter.NewStepRecycleAccumulator(0, int64(len(c.list[i])-1), int64(step))
		if err != nil {
			panic(err)
			return nil
		}

		// 填充cache
		for j := 0; j < rstCount; j++ {
			cache[j][i] = c.list[i][sa.Get()]
			sa.Incr()
		}

	}

	// 转换格式
	var rstList []Result
	for i := 0; i < rstCount; i++ {
		rstItem := Result{}
		rstItem.Items = cache[i]
		rstList = append(rstList, rstItem)
	}
	return rstList
}
