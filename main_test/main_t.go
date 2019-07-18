package main

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"
	"unsafe"
)

func main() {
	//arrsli()
	//stru()
	//run()
	//shunxu()
	//doworlflow()
	//gorou1()
	//chan1()
	//chan2()
	//chan3()
	//chan4()
	//pubAndSub()
	//fastestIsRel()
	Susai()
	//contextf()
	//Susaicontext()
}

func arrsli() {
	//5次循环 0 3 6 7 8
	for i, c := range "\xe4\xb8\x96\xe7\x95\x8cabc" {
		fmt.Println(i, c)
	}
	fmt.Println("----------")
	//9次循环 字节数组
	for i1, c1 := range []byte("世界abc") {
		fmt.Println(i1, c1)
	}
	fmt.Println("----------")
	const s = "\xe4\xb8\x96\xe7\x95\x8cabc"
	for i := 0; i < len(s); i++ {
		fmt.Printf("%d %x \n", i, s[i])
	}
	fmt.Println("----------")
	for i, c := range "\xe4\xb8\x96\xe7\x95\x8cabc" {
		fmt.Printf("%d %s \n", i, string(c))
	}
	fmt.Println("----------")
	fmt.Printf("%#v\n", []rune("世界"))
	fmt.Println("----------")
	fmt.Printf("%#v\n", string([]rune{'世', '界'}))
	fmt.Println("----------")
	forOnstring("12345", func(i int, r rune) {
		fmt.Printf("%d --- %s \n", i, string(r))
	})
	fmt.Println("----------")
	rel := bytestoString([]byte("123"))
	fmt.Println(rel)
	var a = []float64{4.3, 2.1, 5, 7, 2, 1, 88, 1}
	b := a[0:5:7]
	fmt.Println(len(b), "--", cap(b), "--", b)
	sortFloat64V1(a)
	aa := new(A)
	*aa = (A(a))
	sort.Sort(aa)
	fmt.Println(a)
	fmt.Println(aa)
	v := Inc()
	fmt.Println("Inc() rel", v)
}

func stru() {
	//生成方法
	man := new(Man)
	var mf = (*Man).Sex
	mf(man)

	//不关心操作对象的类型
	var mf1 = man.Sex
	mf1()

	var mf2 = func() {
		man.Sex()
		return
	}
	mf2()
}

func doworlflow() {
	worlflow()
	fmt.Println("workflow wait for stop")
	<-stop
}

//for range string 的模拟实现
func forOnstring(s string, forbody func(i int, r rune)) {
	for i := 0; len(s) > 0; {
		r, size := utf8.DecodeRuneInString(s)
		forbody(i, r)
		s = s[size:]
		i = i + size
	}
}

//
func bytestoString(s []byte) string {
	data := make([]byte, len(s))
	for i, c := range s {
		data[i] = c
	}
	var rel string
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&rel))
	hdr.Data = uintptr(unsafe.Pointer(&data[0]))
	hdr.Len = len(s)
	return rel
}

func testcopy() {
	//copy()
}

func sortFloat64V1(a []float64) {
	fmt.Println("a--", a)
	var b []int = (*[1 << 20]int)(unsafe.Pointer(&a[0]))[:len(a):cap(a)]
	fmt.Println("b--", b)
	sort.Ints(b)
	sort.Sort(nil)
	fmt.Println("d_a--", a)
	fmt.Println("d_b--", b)
}

func sortFloat64V2(a []float64) {
	var c []int
	ahdr := (*reflect.SliceHeader)(unsafe.Pointer(&a))
	chdr := (*reflect.SliceHeader)(unsafe.Pointer(&c))
	*chdr = *ahdr
	sort.Ints(c)
}

type A []float64

func (a *A) Len() int {
	return len(*a)
}

func (a *A) Less(i, j int) bool {
	return (*a)[i] < (*a)[j]
}

func (a *A) Swap(i, j int) {
	(*a)[i], (*a)[j] = (*a)[j], (*a)[i]
}

//闭包函数对外部变量的访问方式为引用，几多次访问为一个变量的值
func Inc() (v int) {
	fmt.Println()
	defer func() {
		v++
	}()
	return 42
}

type Man struct {
	Name string
}

func (m *Man) Sex() {
	fmt.Println("NAN")
}

//继承了Man，但不是c++多态的方式，Son.Sex() 实际为Son.Man.Sex(),方法为具体结构体的方法
type Son struct {
	Man
	Age int
}

var total uint64

func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	var i uint64
	for i = 0; i < 10; i++ {
		atomic.AddUint64(&total, i)
	}
}

func run() {
	var wg sync.WaitGroup
	wg.Add(3)
	go worker(&wg)
	go worker(&wg)
	go worker(&wg)
	wg.Wait()
	fmt.Println("toral -- ", total)
}

type singleton struct {
}

var (
	instance    *singleton
	initialized uint32
	mu          sync.Mutex
)

//单例模式
func Instance() *singleton {
	if atomic.LoadUint32(&initialized) == 1 {
		return instance
	}
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		defer atomic.StoreUint32(&initialized, 1)
		instance = &singleton{}
	}
	return instance
}

type Twice struct {
	m    sync.Mutex
	done uint32
}

func (t *Twice) Do(f func()) {
	if atomic.LoadUint32(&t.done) == 2 {
		return
	}
	t.m.Lock()
	defer t.m.Unlock()
	if t.done == 0 || t.done == 1 {
		defer atomic.StoreUint32(&t.done, t.done+1)
		f()
	}
}

var (
	t = Twice{done: 1}
)

func InstanceByTwice() *singleton {
	t.Do(func() {
		instance = &singleton{}
	})
	return instance
}

//atomic.Value
func loadConfig() interface{} {
	return rand.Intn(100)
}

var req chan int

func request() <-chan int {
	t.Do(func() {
		req = make(chan int)
		go func() {
			for {
				select {
				case <-time.Tick(time.Second * 2):
					req <- 1
				}
			}
		}()
	})
	return req
}

var (
	stop = make(chan struct{})
)

func worlflow() {
	var config atomic.Value
	config.Store(loadConfig())
	go func() {
		for {
			time.Sleep(time.Second * 5)
			config.Store(loadConfig())
		}
	}()

	for i := 0; i < 10; i++ {
		go func(i int) {
			for r := range request() {
				c := config.Load()
				fmt.Printf("i:%d -- config %d -- request %d \n", i, c, r)
			}
		}(i)
	}
}

func shunxu() {
	var a string
	var done bool

	setup := func() {
		a = "hello world"
		done = true
	}

	main := func() {
		go setup()
		for !done {
		}
		print(a)
	}
	main()
}

//通一个gorouine满足顺序一致性
func gorou1() {
	main := func() {
		done := make(chan int)
		go func() {
			println("你好，世界")
			done <- 1
		}()
		<-done
	}
	main()
}

func gorou2() {
	main := func() {
		var mu sync.Mutex
		mu.Lock()
		go func() {
			println("你好，世界")
			mu.Unlock()
		}()
		mu.Lock()
	}
	main()
}

func gorou3() {
	var a string
	f := func() {
		print(a)
	}

	//go 语句会在当前Goroutine对应的函数返回前创建新的Goroutine,但go执行的函数
	//f()的执行事件和hello函数的返回的时间是不可排序的，即为并发。
	hello := func() {
		a = "hello,world"
		go f()

	}
	main := func() {
		hello()
	}
	main()
}

func chan1() {
	//无缓存channel上的接受操作总在对应的发送操作完成前发生，并且要保证在两个Goroutinue中执行，否则容易造成死锁。
	//go进程执行的是写入
	var done = make(chan bool)
	var msg string

	aGoroutine := func() {
		msg = "你好，世界"
		done <- true
	}

	main := func() {
		go aGoroutine()
		<-done
		println(msg)
	}
	main()
}

func chan2() {
	//也可以，当写入完成时，msg已经赋值了。注:有风险，对缓存大小敏感，当通道存在缓存大小时不能保证正常打印
	// go进程执行的是接受
	var done = make(chan bool)
	var msg string

	aGoroutine := func() {
		msg = "你好，世界"
		<-done
	}

	main := func() {
		go aGoroutine()
		done <- true
		println(msg)
	}
	main()
}

func chan3() {
	//带缓存通道：对channel的第k次接受操作发生在第k+c次发送操作之前。c为channel的缓存大小.c=0表示为无缓存通道
	var done = make(chan bool, 1)
	var msg string

	aGoroutine := func() {
		msg = "你好，世界"
		<-done
	}

	main := func() {
		go aGoroutine()
		done <- true
		//time.Sleep(time.Second)
		println(msg)
	}
	main()
}

func chan4() {
	var limit = make(chan int, 3)
	var work chan func()
	var g = func() {
		fmt.Println(time.Now())
	}
	t.Do(func() {
		work = make(chan func())
		go func() {
			for {
				select {
				case <-time.Tick(time.Second * 10):
					work <- g
				}
			}
		}()
	})

	main := func() {
		for w := range work {
			go func() {
				limit <- 1
				w()
				<-limit
			}()
		}
		select {}
	}
	main()
}

//too east
func proAndcum() {

}

func pubAndSub() {
	type (
		subscriber chan interface{}

		//true为订阅表示关注，false表示不关注
		topicFunc func(v interface{}) bool
	)

	//发布者对象
	type Publisher struct {
		m           sync.RWMutex
		buffer      int
		timeout     time.Duration
		subscribers map[subscriber]topicFunc
	}

	//构建一个发布者
	Newpublisher := func(publishTimeout time.Duration, buffer int) *Publisher {
		return &Publisher{
			buffer:      buffer,
			timeout:     publishTimeout,
			subscribers: make(map[subscriber]topicFunc),
		}
	}

	closePub := func(p *Publisher) {
		p.m.Lock()
		defer p.m.Unlock()

		for sub := range p.subscribers {
			delete(p.subscribers, sub)
			close(sub)
		}
	}

	//添加一个订阅者，订阅特定主题
	Subscribetopic := func(p *Publisher, topic topicFunc) chan interface{} {
		ch := make(chan interface{}, p.buffer)
		p.m.Lock()
		p.subscribers[ch] = topic
		p.m.Unlock()
		return ch
	}

	//添加一个订阅者，订阅全部主题
	Subscribe := func(p *Publisher) chan interface{} {
		return Subscribetopic(p, nil)
	}

	//退出订阅
	Evict := func(p *Publisher, sub chan interface{}) {
		p.m.Lock()
		defer p.m.Unlock()
		delete(p.subscribers, sub)
		close(sub)
	}
	_ = Evict

	//发送主题，带有超时机制
	sendtopic := func(p *Publisher, sub subscriber, topic topicFunc, v interface{}, wg *sync.WaitGroup) {
		defer wg.Done()
		if topic != nil && !topic(v) {
			//
			return
		}
		select {
		case sub <- v:
		case <-time.After(p.timeout):
		}
	}

	Publish := func(p *Publisher, v interface{}) {
		p.m.Lock()
		defer p.m.Unlock()
		var wg sync.WaitGroup
		for sub, topic := range p.subscribers {
			wg.Add(1)
			go sendtopic(p, sub, topic, v, &wg)
		}
		wg.Wait()
	}

	main := func() {
		wg := sync.WaitGroup{}
		wg.Add(2)
		p := Newpublisher(time.Second, 10)
		defer closePub(p)

		all := Subscribe(p)
		golang := Subscribetopic(p, func(v interface{}) bool {
			if s, ok := v.(string); ok {
				return strings.Contains(s, "golang")
			}
			return false
		})

		Publish(p, "hello world")
		Publish(p, "hello golang")

		go func() {
			for msg := range all {
				fmt.Println("all:=", msg)
			}
			wg.Done()
		}()

		go func() {
			for msg := range golang {
				fmt.Println("golang:= ", msg)
			}
			wg.Done()
		}()

		Evict(p, golang)
		Publish(p, "hello golang")
		Evict(p, all)
		wg.Wait()
	}
	main()
}

func fastestIsRel() {
	var a = make(chan interface{}, 3)
	f1 := func() interface{} {
		n := rand.Intn(20)
		time.Sleep(time.Duration(n) * time.Second)
		return fmt.Sprintf("f1 done -- %d", n)
	}
	f2 := func() interface{} {
		n := rand.Intn(20)
		time.Sleep(time.Duration(n) * time.Second)
		return fmt.Sprintf("f2 done -- %d", n)
	}
	f3 := func() interface{} {
		n := rand.Intn(20)
		time.Sleep(time.Duration(n) * time.Second)
		return fmt.Sprintf("f3 done -- %d", n)
	}

	go func() {
		a <- f1()
	}()
	go func() {
		a <- f2()
	}()
	go func() {
		a <- f3()
	}()
	fmt.Println(<-a)
	fmt.Println(<-a)
	fmt.Println(<-a)
}


func Susai() {
	genNature := func() chan int {
		ch := make(chan int)
		go func() {
			for i := 2; ; i++ {
				ch <- i
			}
			fmt.Println("genNature for end")
			time.Sleep(time.Second)
		}()
		return ch
	}

	primefilter := func(in <-chan int, prime int) chan int {
		out := make(chan int)
		go func() {
			for i := range in {
				if i != 0 && i%prime != 0 {
					out <- i
				}
			}
		}()
		return out
	}

	main := func() {
		ch := genNature()
		for i := 0; i < 100; i++ {
			prime := <-ch
			fmt.Println(i+1, "----", prime)
			ch = primefilter(ch, prime)
		}
	}
	main()
}

func contextf() {
	worker := func(ctx context.Context, wg *sync.WaitGroup) error {
		defer wg.Done()
		for {
			select {
			default:
				time.Sleep(time.Millisecond)
				fmt.Println("hello")
			case <-ctx.Done():
				fmt.Println(ctx.Err())
				return ctx.Err()
			}
		}
	}

	main := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var wg sync.WaitGroup

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go worker(ctx, &wg)
		}
		time.Sleep(20 * time.Second)
		cancel()
		wg.Wait()
	}

	main()
}

func Susaicontext() {
	genNature := func(ctx context.Context) chan int {
		ch := make(chan int)
		go func() {
			for i := 2; ; i++ {
				select {
				case <-ctx.Done():
					return
				case ch <- i:
				}
			}
		}()
		return ch
	}

	primefilter := func(context context.Context, in <-chan int, prime int) chan int {
		out := make(chan int)
		go func() {
			for i := range in {
				if i != 0 && i%prime != 0 {
					select {
					case <-context.Done():
						return
					case out <- i:
					}
				}
			}
		}()
		return out
	}

	main := func() {
		ctx, cancel := context.WithCancel(context.Background())
		ch := genNature(ctx)
		for i := 0; i < 100; i++ {
			prime := <-ch
			fmt.Println(i+1, "----", prime)
			ch = primefilter(ctx, ch, prime)
		}
		cancel()
	}
	main()
}
