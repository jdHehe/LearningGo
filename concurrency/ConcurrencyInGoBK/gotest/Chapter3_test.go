package gotest

import "testing"

func BenchmarkGeneric(b *testing.B){
	done := make(chan  interface{})
	defer close(done)
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{}{
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values{
					select {
					case valueStream <- v:
					case <-done:
						return
					}
				}
			}
		}()
		return valueStream
	}
	take := func(done <- chan interface{}, number int, valueStream <- chan interface{}) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i:=0; i<number; i++{
				select {
				case <-done:
					return
				case takeStream <- <- valueStream:
				}
			}
		}()
		return takeStream
	}
	toString := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan string{
		stringStream := make(chan string)
		go func() {
			defer close(stringStream)
			for v := range valueStream{
				select {
				case <-done:
				case stringStream <- v.(string):
				}
			}
		}()
		return stringStream
	}

	b.ResetTimer()
	for range toString(done, take(done, b.N,repeat(done, "I", "am."))){
	}
}

func BenchmarkTyped(b *testing.B) {
	repeat := func(done <-chan interface{}, values ...string) <-chan string {
		valueStream := make(chan string)
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}
	take := func(
		done <-chan interface{},
		valueStream <-chan string,
		num int,
	) <-chan string {
		takeStream := make(chan string)
		go func() {
			defer close(takeStream)
			for i := num; i > 0 || i == -1; {
				if i != -1 {
					i--
				}
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}
	done := make(chan interface{})
	defer close(done)
	b.ResetTimer()
	for range take(done, repeat(done, "a"), b.N) {
	}
}
