package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

//声明新切片类型
type units []uint32

func (x units) Len() int {
	return len(x)
}

func (x units) Less(i, j int) bool {
	return x[i] < x[j]
}

func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

var errEmpty = errors.New("Hash环没有数据")

//创建结构体保存一致性hash信息
type Consistent struct {
	//hash环， key 为哈希值，值存放节点的信息
	circle map[uint32]string
	//已经排序好的hash切片
	sortedHashes units
	//虚拟节点个数，用来增加hash的平衡性
	VirtualNode int
	//map读写锁
	sync.RWMutex
}

func NewConsistent() *Consistent {
	return &Consistent{
		//初始化变量
		circle:      make(map[uint32]string),
		VirtualNode: 20,
	}
}

//自动生成key值
func (c *Consistent) generateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}

func (c *Consistent) hashkey(key string) uint32 {
	if len(key) < 64 {
		//声明一个数组长度为64
		var srcatch [64]byte
		//拷贝数据到数组
		copy(srcatch[:], key)
		//使用IEEE 多项式返回数据CRC-32校验和
		return crc32.ChecksumIEEE(srcatch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

//更新排序，方便查找
func (c *Consistent) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	//判断切片容量， 是否过大， 若果过大则重置
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.circle) {
		hashes = nil
	}
	//添加hashes
	for k := range c.circle {
		hashes = append(hashes, k)
	}
	//对所有节点hash值进行排序
	sort.Sort(hashes)
	c.sortedHashes = hashes

}

func (c *Consistent) Add(element string) {
	c.Lock()
	defer c.Unlock()
	c.add(element)
}

func (c *Consistent) add(element string) {
	//循环虚拟节点，设置副本
	for i := 0; i < c.VirtualNode; i++ {
		//根据生成的结点添加到hash环中
		c.circle[c.hashkey(c.generateKey(element, i))] = element
	}
	c.updateSortedHashes()
}

func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashkey(c.generateKey(element, i)))

	}
	c.updateSortedHashes()
}

//删除一个结点
func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

//顺时针查找最近节点
func (c *Consistent) serach(key uint32) int {
	//查找算法
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	//使用二分查找算法来搜索指定切片满足条件的最小值
	i := sort.Search(len(c.sortedHashes), f)
	//若果超出范围
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i

}
func (c *Consistent) Get(name string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.circle) == 0 {
		return "", errEmpty
	}
	//计算hash值
	key := c.hashkey(name)
	i := c.serach(key)
	return c.circle[c.sortedHashes[i]], nil
}
