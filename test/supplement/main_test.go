package main

import (
	"fmt"
	"os"
	"testing"
)

func TestAbs(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)
	want := 1
	if got := Abs(-1); got != want {
		t.Fatalf("Abs() = %v, want %v", got, want)
	}
}

func TestAbsWithTable(t *testing.T) {
	type args struct {
		x int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "positive",
			args: args{x: 1},
			want: 1,
		},
		{
			name: "negative",
			args: args{x: -1},
			want: 1,
		},
	}
	for _, tt := range tests {
		teardownTest := setupTest(t)
		defer teardownTest(t) // 错误写法，defer 语句不会在本轮 for 循环结束时被执行
		if got := Abs(tt.args.x); got != tt.want {
			t.Fatalf("Abs() = %v, want %v", got, tt.want)
		}
	}
}

func TestAbsWithTableAndSubtests(t *testing.T) {
	type args struct {
		x int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "positive",
			args: args{x: 1},
			want: 1,
		},
		{
			name: "negative",
			args: args{x: -1},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardownTest := setupTest(t)
			defer teardownTest(t)
			if got := Abs(tt.args.x); got != tt.want {
				t.Fatalf("Abs() = %v, want %v", got, tt.want)
			}
		})
	}
}

// testing.TB is the interface common to T, B, and F.
func setupTest(tb testing.TB) func(tb testing.TB) {
	fmt.Println(">> setup Test")

	return func(tb testing.TB) {
		fmt.Println(">> teardown Test")
	}
}

func setup() {
	fmt.Println("> setup completed")
}

func teardown() {
	fmt.Println("> teardown completed")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
