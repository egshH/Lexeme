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

let input = ""
let n = ""

for (n != val) {
    input = scanln() // getting use input
    n = atoi(input)  // parsing input into INTEGER

    const res = guess(n)
    pritn(res)
}