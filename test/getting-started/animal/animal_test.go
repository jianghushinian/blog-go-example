package animal

import (
	"testing"
)

func TestAnimalDog_shout(t *testing.T) {
	dog := Animal{Name: "dog"}
	got := dog.shout()
	want := "旺！"
	if got != want {
		t.Errorf("got %s; want %s", got, want)
	}
}

func TestAnimalCat_shout(t *testing.T) {
	cat := Animal{Name: "cat"}
	got := cat.shout()
	want := "喵～"
	if got != want {
		t.Errorf("got %s; want %s", got, want)
	}
}
