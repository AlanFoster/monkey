// This file can be run from the root directory of this project with:
//      go run ./main.go --entry-file ./examples/hello-world.monkey
//
puts("Hello" + " " + "world!", 12345, 1 + 2 + 3)

let add = fn(x, y) { x + y; }
puts("Function calls: ", add(1, 2))

let createAdder = fn(x) { fn(y) { x + y } }
let addFive = createAdder(5)
puts("Closure Example: ", addFive(8))

let array = [5, 4, 3, 2, 1]
puts("Array example: ", array)

puts("🎉 🎉 🎉 le fin! 🎉 🎉 🎉")