package main

type Msg struct {
	Mensagem string `json:"mensagem"`
}

var Message = Msg{"Esta rodando"}

func main() {
	Loading()
}
