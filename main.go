package main

import (
    "bytes"
    "encoding/csv"
    "errors"
    "flag"
    "io"
    "io/ioutil"
    "regexp"
    "sort"
    "strings"
    "os"
    "fmt"
)

type SubstsTree map[string][]string
type SubstsMap map[string]string

var inFile string
var outFile string
var partNumberFromColumn int
var partNumberToColumn int

func init() {
    flag.StringVar(&inFile, "src", "", "source/input file name")
    flag.StringVar(&outFile, "dst", "/dev/stdout", "destination/output file name (default /dev/stdout)")
    flag.IntVar(&partNumberFromColumn, "a", 0, "part number A (from) column number")
    flag.IntVar(&partNumberToColumn, "b", 0, "part number B (to) column number")
    flag.Parse()
}

func main() {
    if len(inFile) == 0 {
        checkError(errors.New("flag --src is missing"))
    }

    if len(outFile) == 0 {
        checkError(errors.New("flag --out is missing"))
    }

    if partNumberFromColumn == 0 {
        checkError(errors.New("flag --a is missing"))
    }

    if partNumberToColumn == 0 {
        checkError(errors.New("flag --b is missing"))
    }

    if (partNumberFromColumn == 0 || partNumberToColumn == 0) || (partNumberFromColumn == partNumberToColumn) {
        checkError(errors.New("--a and --b shoud be greater than 0 and not be equal"))
    } else {
        partNumberFromColumn--
        partNumberToColumn--
    }

    substitutes := make(SubstsTree, 1024)
    substitutesCircle := make(SubstsMap, 1024)
    substitutesHistory := make(SubstsMap, 4)

    reg, err := regexp.Compile("[^0-9A-Z]+")
    checkError(err)

    data, err := ioutil.ReadFile(inFile)
    checkError(err)

    r := csv.NewReader(bytes.NewReader(data))
    r.Comma = '\t'
    r.Comment = '#'
    r.FieldsPerRecord = 4

    for {
        record, err := r.Read()
        if err == io.EOF {
            break
        } else {
            checkError(err)
        }

        partNumberFrom := strings.ToUpper(record[partNumberFromColumn])
        partNumberFrom = reg.ReplaceAllString(partNumberFrom, "")

        partNumberTo := strings.ToUpper(record[partNumberToColumn])
        partNumberTo = reg.ReplaceAllString(partNumberTo, "")

        if len(partNumberFrom) > 0 && len(partNumberTo) > 0 {
            substitutes[partNumberFrom] = append(substitutes[partNumberFrom], partNumberTo)
        }
    }

    for key, substs := range substitutes {
        result := checkCircle(key, substs, &substitutes, &substitutesHistory)
        if result == false {
            substitutesCircle[key] = key
        }

        substitutesHistory = make(SubstsMap, 4)
    }

    if len(substitutesCircle) == 0 {
        return
    }

    substsKeys := make([]string, 16)

    for key := range substitutesCircle {
        substsKeys = append(substsKeys, key)
    }

    sort.Strings(substsKeys)

    result := ""
    for k := 0; k < len(substsKeys); k++ {
        if substsKeys[k] != "" {
            result += substsKeys[k] + "\n"
        }
    }

    err = ioutil.WriteFile(outFile, []byte(result), 0644)
    checkError(err)
}

func checkCircle(key string, sliceOfSubstitutes []string, substitutes *SubstsTree, checked *SubstsMap) bool {
    for _, el := range sliceOfSubstitutes {
        if _, exist := (*checked)[el] ; exist {
            continue
        } else {
            (*checked)[el] = el
        }

        if el == key {
            return false
        } else {
            (*checked)[el] = el
            if value, ok := (*substitutes)[el]; ok {
                if checkCircle(key, value, substitutes, checked) == false {
                    return false
                }
            }
        }
    }

    return true
}

func checkError(e error) {
    if e != nil {
        fmt.Printf("error: %s\n", e)
        os.Exit(1)
    }
}
