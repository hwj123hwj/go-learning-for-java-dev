package main

import "fmt"

type User struct {
	Name  string
	Age   int
	Email string
}

type Admin struct {
	User
	Role string
}

func (u User) SayHello() {
	fmt.Printf("%s (Age: %d) 说: Hello!\n", u.Name, u.Age)
}

func (u *User) SetAge(age int) {
	u.Age = age
}

func main() {
	fmt.Println("=== 结构体与方法 (对比Java class) ===")
	fmt.Println("")

	user := User{Name: "Alice", Age: 25, Email: "alice@example.com"}
	fmt.Printf("创建User: %+v\n", user)

	user.SayHello()
	user.SetAge(26)
	fmt.Printf("调用SetAge(26)后: %+v\n", user)

	admin := Admin{
		User: User{Name: "Bob", Age: 30, Email: "bob@example.com"},
		Role: "admin",
	}
	fmt.Printf("\nAdmin继承User: %+v\n", admin)
	admin.SayHello()

	user2 := new(User)
	user2.Name = "Charlie"
	user2.Age = 28
	user2.Email = "charlie@example.com"
	fmt.Printf("\n使用new创建: %+v\n", user2)
}
