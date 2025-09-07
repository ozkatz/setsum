# go-setsum

This is a Go port of [setsum](https://github.com/rescrv/blue/tree/main/setsum): An associative and commutative checksum that offers standard multi-set operations over strings in constant time.

From the original README:

> setsum is an unordered checksum that operates on sets of data.  Where most
> checksum algorithms process data in a streaming fashion and produce a checksum
> that's unique for each stream, setsum processes a stream of discrete elements
> and produces a checksum that is unique to those elements, regardless of the
> order in which they occur in the stream.


## Example Usage

```go
import (
    "bytes"
    "log"

    "github.com/ozkatz/setsum"
)


func SetsumExample() {
    sum1 := setsum.Default()
    sum1.Insert([]byte("A"))
    sum1.Insert([]byte("B"))

    sum2 := setsum.Default()
    sum2.Insert([]byte("B"))
    sum2.Insert([]byte("A")

    if !bytes.Equal(sum1.Digest(), sum2.Digest()) {
        log.Fatal("both should be equal")
    }

    // or simply,
    if !sum1.Equals(sum2) {
        log.Fatal("both should be equal")
    }
}
```

See the [tests](https://github.com/ozkatz/setsum/tree/main/setsum_test.go) for more usage examples.


## License

This project is licensed under the Apache License, Version 2.0.

It is a Go port of the Rust library [github.com/rescrv/blue/tree/main/setsum](https://github.com/rescrv/blue/tree/main/setsum),
originally developed by Dropbox.
