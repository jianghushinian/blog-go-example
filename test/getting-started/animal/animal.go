package animal

type Animal struct {
	Name string
}

func (a Animal) shout() string {
	if a.Name == "dog" {
		return "旺！"
	}
	if a.Name == "cat" {
		return "喵～"
	}
	return "吼～"
}
