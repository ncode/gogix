package util

import (
    "log"
)

func Checkp(err error) {
    if err != nil {
        panic(err)
    }
}

func Checkl(err error) {
    if err != nil {
        log.Fatal(err)
    }
}