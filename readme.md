# Lexeme lang

Lexeme is a little interpreted programming language, made by following on book "writing interpreter in go" by Thorston Ball.

## Installation
```bash
git clone https://github.com/egshH/Lexeme.git
cd Lexeme/cmd/lexeme
go install # or with go build
```

## Usage
### Examples

You can fild a few examples of lexeme scripts int `./example/` folder.

### Getting started

Let's assume that you already installed a `lexeme` interpreter. To run your first script, you should provide a path to script file with `.lex` ext.
```shell
lexeme foo/bar.lex
```
In this example we running `bar.lex` file in `foo` dir.
Let's write our first script!. Instead of creating file we can use REPL. run `lexeme` and you will start it.

Anyway, let's write hello world!
```js
print("hello world!") // 'hello world!'
```
And this is it! Our first script!
To quit REPL press `<ctrl+c>`

### Syntax
#### Variables
```js
const _my_constant = 5; // this value can't be changed
let _my_var = "hello";  //this one - can

_my_constant = 6; // => error
_my_var = 123;   // ok, now _my_var is an integer
```

##### Strings
```js
let s = "new"
s = s + " text"
let f = "another text"

let line = s + " " + f
print(line) // 'new text another text'

print(s[0]) // symbol - 'n'
```

##### Arrays
```js
let arr = [1, 2, 3]
print(arr[0]) // 1

arr = push(arr, "line")
print(arr) // [1, 2, 3, "line"], yes, arrays can contain different type of values
```

##### Hashmaps
```js
let m = {"name": "foo", "age": 21}
print(m) // {"name": "foo", "age": 21}

m["job"] = "jobless"
print(m) //{"name": "foo", "age": 21, "job": "jobless"}

m["job"] = "qa engineer"
print(m) //{"name": "foo", "age": 21, "job": "qa engineer"}
```

##### Boolean
```js
// boolean
let t = true
let f = false

let n = t || f // true
let m = t && f // false
let o = !t // false
let p = !f // true
```

#### Functions
```js
let add = fn(x,y) {
	return x+y;
}  // you can also type only x+y, without return keyword

let res = add(5, 10)
print(res) // 15
```

#### Control flow
```js
let dev = fn(x,y) {
	if (y==0) {
		return;
	} // if you wanna return from function, you should use return keyword with ';' semicolon symbol. This way function will ruturn null
	
	x/y; // just as I said earlier, this will return x/y, even without return keyword
}

const foo = dev(20, 5)
const bar = dev(20, 0)

print(foo, bar) // 4, null
```

Les't create guess.lex file (you can see it in example dir) and write a little guess game
```js
// guess.lex
const val = 55;
const guess = fn(x) {
	if (x == val) {
		return "you guessed it!"
	} else {
		
		if (x < val) {
			return "too low"
		} else {
			return "too high"
		}
	}
}

let res = guess(123)
print(res) // => too high
```
> [!warning]
> lexeme doesn't support `else if (cond)` statement

#### Loops
let's continue on working our 'guess' game
```js
// guess.lex
// [...]

let input = ""
let n = ""
for (n != val) {
	input = scanln() // getting user input
	n = atoi(input)  // parsing input to INTEGER
	
	const res = guess(n)
	print(res)
}
```

In lexeme there are no 'while' loops, only 'for' loops
```js
for (cond) {}
```

If you want to run your loop forever, you can do this:
```js
for {
	// this will run forever
}
```

If you need to break out of the loop or skip one iteration you can use `break` and `continue` keywords.

#### Error Handling
You know what? F@CK try-catch! But here it is:
```js
try{
	FuncWithError()
} catch(e) {
	print(e)
}
```

Errors are different types. They are similar to strings, but errors.
You can provide any name in brackets and error message will be assigned to it. This error is constant, you **can't** reassign it!

Also, you can ignore `catch` block:
```js
try {
	const data = read("unknown file")
	print(data)
} // Error will be ignored
```