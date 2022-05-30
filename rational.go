package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Floater64 interface {
	// Converts a value to an equivalent float64.
	toFloat64() float64
}

type Rationalizer interface {

	// 5. Rationalizers implement the standard Stringer interface.
	fmt.Stringer

	// 6. Rationalizers implement the Floater64 interface.
	Floater64

	// 2. Returns the numerator.
	Numerator() int

	// 3. Returns the denominator.
	Denominator() int

	// 4. Returns the numerator, denominator.
	Split() (int, int)

	// 7. Returns true iff this value equals other.
	Equal(other Rationalizer) bool

	// 8. Returns true iff this value is less than other.
	LessThan(other Rationalizer) bool

	// 9. Returns true iff the value equal an integer.
	IsInt() bool

	// 10. Returns the sum of this value with other.
	Add(other Rationalizer) Rationalizer

	// 11. Returns the product of this value with other.
	Multiply(other Rationalizer) Rationalizer

	// 12. Returns the quotient of this value with other. The error is nil
	// if its is successful, and a non-nil if it cannot be divided.
	Divide(other Rationalizer) (Rationalizer, error)

	// 13. Returns the reciprocal. The error is nil if it is successful,
	// and non-nil if it cannot be inverted.
	Invert() (Rationalizer, error)

	// 14. Returns an equal value in lowest terms.
	ToLowestTerms() Rationalizer
} // Rationalizer interface

type Rational struct {
	numerator, denominator int
}

// 1. Make a rational
func makeRational(n, d int) (Rational, error) {
	if d == 0 {
		return Rational{0, 0}, errors.New("cannot make a rational number with denominator equals to 0 ")
	}

	return Rational{n, d}, nil
}

// 2. Get the numerator
func (r Rational) Numerator() int {
	return r.numerator
}

// 3. Get the denominator
func (r Rational) Denominator() int {
	return r.denominator
}

// 4. Get the numerator and denominator as a pair
func (r Rational) Split() (int, int) {
	return (r.numerator), (r.denominator)
}

// 5. Convert to a string
func (r Rational) String() string {
	return fmt.Sprintf("%v/%v", r.numerator, r.denominator)
}

// 6. Convert the fraction to a floating point value
func (r Rational) toFloat64() float64 {
	return float64(r.numerator) / float64(r.denominator)
}

// 7. Test for equality
func (r Rational) Equal(other Rationalizer) bool {

	currentR := Rational{r.ToLowestTerms().Numerator(), r.ToLowestTerms().Denominator()}
	otherR := Rational{other.ToLowestTerms().Numerator(), other.ToLowestTerms().Denominator()}

	return currentR.numerator*otherR.Denominator() == currentR.denominator*otherR.Numerator()
}

// 8. Test for order: less than
func (r Rational) LessThan(other Rationalizer) bool {

	if r.denominator < 0 {
		r.numerator *= -1
		r.denominator *= -1
	}
	otherNum := other.ToLowestTerms().Numerator()
	otherDen := other.ToLowestTerms().Denominator()

	if otherDen < 0 {
		otherDen *= -1
		otherNum *= -1
	}

	if otherDen == 0 {
		return r.toFloat64() < 0
	}
	return r.numerator*otherDen < otherNum*r.denominator
}

// 9. Tests if an integer
func (r Rational) IsInt() bool {
	if r.denominator == 0 {
		return false
	}
	return r.numerator%r.denominator == 0
}

// 10. Add
func (r Rational) Add(other Rationalizer) Rationalizer {
	newNum := r.numerator*other.Denominator() + other.Numerator()*r.denominator

	if newNum != 0 {
		newDen := r.denominator * other.Denominator()

		if newNum == newDen {
			return Rational{1, 1}
		} else {
			if newDen < 0 {
				newNum *= -1
				newDen *= -1
			}
			return Rational{newNum, newDen}.ToLowestTerms()
		}
	}
	return Rational{0, 1}
}

// 11. Multiply
func (r Rational) Multiply(other Rationalizer) Rationalizer {
	newNum := r.numerator * other.Numerator()
	newDen := r.denominator * other.Denominator()

	if newNum != 0 {
		if newNum == newDen {
			return Rational{1, 1}
		} else {
			if newDen < 0 {
				newNum *= -1
				newDen *= -1
			}
			return Rational{newNum, newDen}.ToLowestTerms()
		}
	}
	return Rational{0, newDen}
}

// 12. Divide
func (r Rational) Divide(other Rationalizer) (Rationalizer, error) {
	newDen := r.denominator * other.Numerator()
	if r.denominator != 0 && other.Numerator() != 0 && other.Denominator() != 0 {
		newNum := r.numerator * other.Denominator()

		if newNum == newDen {
			return Rational{1, 1}, nil
		}

		if newDen < 0 {
			newNum *= -1
			newDen *= -1
		}
		returnVal, err := makeRational(newNum, newDen)
		return returnVal.ToLowestTerms(), err

	} else if r.denominator == 0 || other.Numerator() == 0 {
		return r, errors.New(" is not divisible, because either denominator of first rational number is 0 or numerator of second rational number is 0")
	}
	return makeRational(0, newDen)
}

// 13. Invert
func (r Rational) Invert() (Rationalizer, error) {
	if r.numerator == 0 {
		return r, errors.New("Numerator is 0, so it has no inverse")
	}
	return Rational{r.denominator, r.numerator}, nil
}

// 14. Returns an equal value in lowest terms
func (r Rational) ToLowestTerms() Rationalizer {
	if r.denominator < 0 {
		r.numerator *= -1
		r.denominator *= -1
	}

	curNumerator := r.numerator
	curDenominator := r.denominator

	if curNumerator < 0 && curDenominator > 0 {
		curNumerator *= -1
	} else if curDenominator < 0 && curNumerator >= 0 {
		curDenominator *= -1
	}

	n := curNumerator
	d := curDenominator

	var gcd int
	for i := 1; i <= n && i <= d; i++ {
		if n%i == 0 && d%i == 0 {
			gcd = i
		}
	}

	if gcd > 1 {
		n = r.numerator / gcd
		d = r.denominator / gcd

		return Rational{n, d}
	}

	return Rational{r.numerator, r.denominator}
}

// 15. Harmonic sum
func makeHarmonicSum(n int) Rationalizer {
	var sum Rationalizer
	sum = Rational{1, 1}
	for i := 2; i <= n; i++ {
		newRational := Rational{1, i}
		sum = sum.Add(newRational)
	}
	return sum
}

// Insertion sort for int
func sortInt(lst []int) []int {
	newList := make([]int, len(lst)) // Creating new list for the return value

	copy(newList, lst) // Copies the elements of the parameter list to the new list

	for i := 0; i < len(newList); i++ {
		for j := i; j > 0 && newList[j-1] > newList[j]; j-- {
			newList[j], newList[j-1] = newList[j-1], newList[j] // if the value at next index is smaller than the current one, swap
		}
	}
	return newList
}

//	Insertion sort for string
func sortString(lst []string) []string {
	newList := make([]string, len(lst))

	copy(newList, lst) // Copies the elements of the parameter list to the new list

	for i := 0; i < len(newList); i++ {
		for j := i; j > 0 && newList[j-1] > newList[j]; j-- {
			newList[j], newList[j-1] = newList[j-1], newList[j] // if the value at next index is smaller than the current one, swap
		}
	}
	return newList
}

// Insertion sort for rational numbers
func sortRational(lst []Rational) []Rational {
	newList := make([]Rational, len(lst))

	copy(newList, lst) // Copies the elements of the parameter list to the new list

	for i := 0; i < len(newList); i++ {
		for j := i; j > 0 && newList[j-1].LessThan(newList[j]); j-- {
			newList[j], newList[j-1] = newList[j-1], newList[j] // if the value at next index is smaller than the current one, swap
		}
	}
	return newList
}

// To generate a list with n random strings

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Testing
func main() {
	// Code that generates lists and timings for Sorting Performance Experiment

	n := 10000 // List size

	// Integer list test
	a := make([]int, n)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		a[i] = rand.Intn(100)
	}

	start := time.Now()
	sortInt(a)
	duration := time.Since(start)

	fmt.Println("Sorting for the list of", n, "integers took ", (duration.Microseconds()), "microseconds to execute!")

	// String list test
	b := make([]string, n)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		b[i] = randSeq(10)
	}

	start = time.Now()
	sortString(b)
	duration = time.Since(start)
	fmt.Println("Sorting for the list of", n, "strings took ", (duration.Microseconds()), "microseconds to execute!")

	// Rational number list test
	c := make([]Rational, n)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		num := rand.Intn(50)
		den := rand.Intn(50-1) + 1
		c[i] = Rational{num, den}
	}

	start = time.Now()
	sortRational(c)
	duration = time.Since(start)
	fmt.Println("Sorting for the list of", n, "rational numbers took ", (duration.Microseconds()), "microseconds to execute!")

	// --------------------------------------------------------------------------- //

	// Rational type function tests

	r, err := makeRational(-1, 5)

	if err != nil {
		fmt.Printf("error: %s", err.Error())
	} else {
		fmt.Println("Rational Number: ", r)
		fmt.Println("Numerator: ", r.Numerator(), "Denominator: ", r.Denominator())

		var num, denom = r.Split()
		fmt.Println("Pair: ", num, denom)
	}

	harmonicSum := makeHarmonicSum(6)
	fmt.Println("Rational Number: ", harmonicSum)

	inverted, err := harmonicSum.Invert()
	fmt.Println(inverted, err)

	fmt.Println(inverted.Equal(harmonicSum))

	newRational := Rational{6, 0}

	fmt.Println(inverted.Equal(newRational))
	fmt.Println(inverted.toFloat64())
	fmt.Println(newRational.LessThan(r))
	fmt.Println(newRational.IsInt())

	testMulti := r.Multiply(newRational)

	fmt.Println(testMulti)

	testDiv, err := r.Divide(newRational)

	fmt.Println(testDiv, err)
}
