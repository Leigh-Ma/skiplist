package skiplist

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

type Int struct {
	v int
	r int
	id int64
}

func NewInt(i int) *Int {
	return &Int{
		v: i,
		id: int64(i),
	}
}

func (i Int) Less(other interface{}) bool {
	return i.v < other.(*Int).v
}

func (i Int) Rank() int {
	return i.r
}

func (i Int) ObjectId() interface{} {
	return i.id
}

func (i *Int) SetRank(r int) int {
	i.r = r
	return i.r
}

func TestInt(t *testing.T) {
	sl := New()
	defer output(sl)
	if sl.Len() != 0 || sl.Front() != nil && sl.Back() != nil {
		t.Fatal()
	}

	testData := []*Int{NewInt(1), NewInt(2), NewInt(3)}

	sl.Insert(testData[0])
	if sl.Len() != 1 || sl.Front().Value.(*Int) != testData[0] || sl.Back().Value.(*Int) != testData[0] {
		t.Fatal()
	}

	sl.Insert(testData[2])
	if sl.Len() != 2 || sl.Front().Value.(*Int) != testData[0] || sl.Back().Value.(*Int) != testData[2] {
		t.Fatal()
	}

	sl.Insert(testData[1])
	if sl.Len() != 3 || sl.Front().Value.(*Int) != testData[0] || sl.Back().Value.(*Int) != testData[2] {
		t.Fatal()
	}

	sl.Insert(NewInt(-999))
	sl.Insert(NewInt(-888))
	sl.Insert(NewInt(888))
	sl.Insert(NewInt(999))
	sl.Insert(NewInt(1000))

	expect := []*Int{NewInt(-999), NewInt(-888), NewInt(1), NewInt(2), NewInt(3), NewInt(888), NewInt(999), NewInt(1000)}
	ret := make([]*Int, 0)

	for e := sl.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*Int))
	}
	for i := 0; i < len(ret); i++ {
		if ret[i].v != expect[i].v {
			t.Fatal()
		}
	}

	e := sl.Find(NewInt(2))
	if e == nil || e.Value.(*Int).v != 2 {
		t.Fatal()
	}

	ret = make([]*Int, 0)
	for ; e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*Int))
	}
	for i := 0; i < len(ret); i++ {
		if ret[i].v != expect[i+3].v {
			t.Fatal()
		}
	}

	sl.Remove(sl.Find(NewInt(2)))
	sl.Delete(NewInt(888))
	sl.Delete(NewInt(1000))

	expect = []*Int{NewInt(-999), NewInt(-888), NewInt(1), NewInt(3), NewInt(999)}
	ret = make([]*Int, 0)

	for e := sl.Back(); e != nil; e = e.Prev() {
		ret = append(ret, e.Value.(*Int))
	}

	for i := 0; i < len(ret); i++ {
		if ret[i].v != expect[len(ret)-i-1].v {
			t.Fatal()
		}
	}

	if sl.Front().Value.(*Int).v != -999 {
		t.Fatal()
	}

	sl.Remove(sl.Front())
	if sl.Front().Value.(*Int).v != -888 || sl.Back().Value.(*Int).v != 999 {
		t.Fatal()
	}

	sl.Remove(sl.Back())
	if sl.Front().Value.(*Int).v != -888 || sl.Back().Value.(*Int).v != 3 {
		t.Fatal()
	}

	if e = sl.Insert(NewInt(2)); e.Value.(*Int).v != 2 {
		t.Fatal()
	}
	sl.Delete(NewInt(-888))

	if r := sl.Delete(NewInt(123)); r != nil {
		t.Fatal()
	}

	if sl.Len() != 3 {
		t.Fatal()
	}

	sl.Insert(NewInt(2))
	sl.Insert(NewInt(2))
	sl.Insert(NewInt(1))

	if e = sl.Find(NewInt(2)); e == nil {
		t.Fatal()
	}

	expect = []*Int{NewInt(2), NewInt(2), NewInt(2), NewInt(3)}
	ret = make([]*Int, 0)
	for ; e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*Int))
	}
	for i := 0; i < len(ret); i++ {
		if ret[i].v != expect[i].v {
			t.Fatal()
		}
	}

	sl2 := sl.Init()
	if sl.Len() != 0 || sl.Front() != nil || sl.Back() != nil ||
		sl2.Len() != 0 || sl2.Front() != nil || sl2.Back() != nil {
		t.Fatal()
	}

	// for i := 0; i < 100; i++ {
	// 	sl.Insert(Int(rand.Intn(200)))
	// }
	// output(sl)
}

func TestRank(t *testing.T) {
	sl := New()
	defer output(sl)
	for i := 1; i <= 10; i++ {
		sl.Insert(NewInt(i))
	}

	for i := 1; i <= 10; i++ {
		if sl.GetRankFast(NewInt(i)) != i {
			t.Fatal(i, sl.GetElementByRank(i).Value)
		}
	}

	for i := 1; i <= 10; i++ {
		if sl.GetElementByRank(i).Value.(*Int).v != i{
			t.Fatal()
		}
	}

	if sl.GetRankFast(NewInt(0)) != 0 || sl.GetRank(NewInt(11)) != 0 {
		t.Fatal()
	}

	if sl.GetElementByRank(11) != nil || sl.GetElementByRank(12) != nil {
		t.Fatal()
	}

	expect := []*Int{NewInt(7), NewInt(8), NewInt(9), NewInt(10)}
	for e, i := sl.GetElementByRank(7), 0; e != nil; e, i = e.Next(), i+1 {
		if e.Value.(*Int).v != expect[i].v {
			t.Fatal()
		}
	}

	sl = sl.Init()
	mark := make(map[int]bool)
	ss := make([]int, 0)

	for i := 1; i <= 100000; i++ {
		x := rand.Int()
		if !mark[x] {
			mark[x] = true
			sl.Insert(NewInt(x))
			ss = append(ss, x)
		}
	}
	sort.Ints(ss)

	for i := 0; i < len(ss); i++ {
		if sl.GetElementByRank(i+1).Value.(*Int).v != NewInt(ss[i]).v || sl.GetRank(NewInt(ss[i])) != i+1 {
			t.Fatal()
		}
	}

	// output(sl)
}

func BenchmarkIntInsertOrder(b *testing.B) {
	b.StopTimer()
	sl := New()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.Insert(NewInt(i))
	}
}

func BenchmarkIntInsertRandom(b *testing.B) {
	b.StopTimer()
	sl := New()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.Insert(NewInt(rand.Int()))
	}
}

func BenchmarkIntDeleteOrder(b *testing.B) {
	b.StopTimer()
	sl := New()
	for i := 0; i < 1000000; i++ {
		sl.Insert(NewInt(i))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.Delete(NewInt(i))
	}
}

func BenchmarkIntDeleteRandome(b *testing.B) {
	b.StopTimer()
	sl := New()
	for i := 0; i < 1000000; i++ {
		sl.Insert(NewInt(rand.Int()))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.Delete(NewInt(rand.Int()))
	}
}

func BenchmarkIntFindOrder(b *testing.B) {
	b.StopTimer()
	sl := New()
	for i := 0; i < 1000000; i++ {
		sl.Insert(NewInt(i))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.Find(NewInt(i))
	}
}

func BenchmarkIntFindRandom(b *testing.B) {
	b.StopTimer()
	sl := New()
	for i := 0; i < 1000000; i++ {
		sl.Insert(NewInt(rand.Int()))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.Find(NewInt(rand.Int()))
	}
}

func BenchmarkIntRankOrder(b *testing.B) {
	b.StopTimer()
	sl := New()
	for i := 0; i < 1000000; i++ {
		sl.Insert(NewInt(i))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.GetRank(NewInt(i))
	}
}

func BenchmarkIntRankRandom(b *testing.B) {
	b.StopTimer()
	sl := New()
	for i := 0; i < 1000000; i++ {
		sl.Insert(NewInt(rand.Int()))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.GetRank(NewInt(rand.Int()))
	}
}

func output(sl *SkipList) {
	s := time.Now().UnixNano()
	sl.refreshDirtyRank(sl.Len() + 1)
	fmt.Printf("refresh %d elements rank info use %d ns\n", sl.Len(), time.Now().UnixNano() - s)
	/*for x, r := sl.header.level[0].forward, 1; x != nil; x, r=x.Next(),r+1 {
		fmt.Printf("[%2d]  %v\n", r, x.Value)

		for l := 0; l < len(x.level); l++ {
			fmt.Print("\t")
			if x.level[l] != nil {
				for j:=0; j <= l; j++ {
					fmt.Print("--> ")
				}
				fmt.Printf("skip %d\n", x.level[l].span)
			}

		}

	}
	*/

	fmt.Printf("====================\n")
}
