package main

func Foo() vec2 {
	x := Float() * Int()
	y := Int() * Float()
	return vec2(x, y)
}

func Foo2() vec2 {
	var x = Float() * Int()
	var y = Int() * Float()
	return vec2(x, y)
}

func Float() float {
	return 1.0
}

func Int() int {
	return 1
}
