package main

import (
    "bytes"
    "encoding/csv"
    "errors"
    "flag"
    "io"
    "io/ioutil"
    "log"
    "regexp"
    "sort"
    "strings"
)

var inFile string
var outFile string
var partNumberFromColumn int
var partNumberToColumn int

func init() {
    flag.StringVar(&inFile, "src", "", "source (input) file name")
    flag.StringVar(&outFile, "dst", "", "destination (output) file name")
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

    substitutes := make(map[string][]string, 1024)
    substitutesCircle := make(map[string]string, 1024)

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
        result := checkCircle(0, key, substs, &substitutes)
        if result == false {
            substitutesCircle[key] = key
        }
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

// check circle substitutes
func checkCircle(idx int, key string, sliceOfSubstitutes []string, substitutes *map[string][]string) bool {
    if idx > 10 {
        return true
    } else {
        idx++
    }

    for _, el := range sliceOfSubstitutes {
        if el == key {
            return false
        } else {
            if value, ok := (*substitutes)[el]; ok {
                return checkCircle(idx, key, value, substitutes)
            }
        }
    }

    return true
}

// check for error
func checkError(e error) {
    if e != nil {
        log.Fatal(e)
    }
}
