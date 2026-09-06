package main

import (
	"fmt"
	"sync"
)

type Singleton struct {
	Name string
}

var instance *Singleton
var once sync.Once

func GetInstance() *Singleton {
	once.Do(func() {
		fmt.Println("初始化Singleton...")
		instance = &Singleton{Name: "我是单例"}
	})
	return instance
}

func main() {
	fmt.Println("=== sync.Once ===")
	fmt.Println()

	for i := 0; i < 5; i++ {
		s := GetInstance()
		fmt.Printf("第%d次获取: %p, name=%s\n", i+1, s, s.Name)
	}

	fmt.Println()
	fmt.Println("结果：只初始化一次，所有实例地址相同")
	fmt.Println()
	fmt.Println("--- Java vs Go 对比 ---")
	fmt.Println()
	fmt.Println("Java双检锁单例:")
	fmt.Println(`  private static volatile Singleton instance;
  public static Singleton getInstance() {
      if (instance == null) {
          synchronized (Singleton.class) {
              if (instance == null) {
                  instance = new Singleton();
              }
          }
      }
      return instance;
  }`)
	fmt.Println()
	fmt.Println("Go sync.Once:")
	fmt.Println(`  var instance *Singleton
  var once sync.Once
  func GetInstance() *Singleton {
      once.Do(func() {
          instance = &Singleton{}
      })
      return instance
  }`)
	fmt.Println()
	fmt.Println("Go惯用法:")
	fmt.Println("  - sync.Once是Go单例的标准做法")
	fmt.Println("  - 代码简洁，线程安全")
	fmt.Println("  - 无论调用多少次Do，函数只执行一次")
}
