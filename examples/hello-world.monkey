// This file can be run from the root directory of this project with:
//      go run ./main.go --entry-file ./examples/hello-world.monkey
//
puts("Hello" + " " + "world!", 12345, 1 + 2 + 3)

let add = fn(x, y) { x + y; }
puts("Function calls: ", add(1, 2))

// Closure Examples
let createAdder = fn(x) { fn(y) { x + y } }
let addFive = createAdder(5)
puts("Closure Example: ", addFive(8))

// Array examples
let reduce = fn(array, initial, reducer) {
  let iter = fn(array, acc, reducer) {
    if (len(array) == 0) {
        acc
    } else {
        iter(rest(array), reducer(acc, first(array)), reducer)
    }
  }

  iter(array, initial, reducer)
}

let map = fn(array, f) {
  reduce(array, [], fn(acc, next) { push(acc, f(next)) })
}

let array = [1, 2, 3, 4, 5]
puts("Starting array: ", array)

let succ = fn(x) { x + 1 };
puts("Mapping over the array: ", map(array, succ)) // Arrays in monkey are immutable

let sum = fn(array) { reduce(array, 0, fn(acc, next) { acc + next }) }
puts("Sum of original list:", sum(array))

puts("🎉 🎉 🎉 🎉 🎉 🎉")
