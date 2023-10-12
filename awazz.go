package main

// "github.com/SohelAhmedJoni/Awazz-Backend/internal/model"
import "github.com/SohelAhmedJoni/Awazz-Backend/internal/model"

func main() {
	saria := model.Account{}
	saria.Email = "saria"
	saria.Password = "rumpush"
	saria.Bio = "here's is my thought story idea"
	println(saria.Email)
	println(saria.Password)
	println(saria.Bio)
}
