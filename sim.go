package main

import (
  "fmt"
  "os"
  "bufio"
  "strconv"
  "strings"
  "math"
  "log"
)

func main() {
  if len(os.Args) < 4 {
    fmt.Println("Missing parameter. Please check to make sure all arguments are present.")
    return
  }

  M,err := strconv.Atoi(os.Args[1])
  N,err := strconv.Atoi(os.Args[2])

  file, err := os.Open(os.Args[3])

  if err != nil {
    fmt.Println("Can't read file: ", os.Args[3])
    return
  }
  defer file.Close();

  branch_num := 0
  miss := 0
  branchVal := -1

  a := make([]int, int (math.Pow(2, float64(M))))

  for i := range a {
    a[i] = 2
  }

  GBH := 0
  //var gtable data

  scanner := bufio.NewScanner(file)
  //iterate through trace file
  for scanner.Scan() {
    var x = strings.Fields(scanner.Text())

    pcAddress, err := strconv.ParseUint(x[0], 16, 64) //uint64
    branchResult := x[1] //string

    //check error during conversion
    if err != nil {
      fmt.Println(err)
      return
    }

    if branchResult == "t" {
      branchVal = 1
    } else if branchResult == "n" {
      branchVal = 0
    } else {
      panic("Unexpected branch result.")
    }

    //Divides pcAddress by 4 or removes 2 LSB
    pcTrimmed := pcAddress>>2
    //Moves M-N 0's to GBH
    M_GBH := GBH<<(uint64(M)-uint64(N))
    //Takes the M bits of the trimmed PC
    M_PC := pcTrimmed&((1<<uint64(M))-1)
    //Gets i from XOR of M_GBH and M_PC
    i := uint64 (M_GBH)^M_PC
    //Gets state at index i of array a
    state := a[i]
    //Taken branch
    if (branchVal == 1) {
      //Check if it was Miss Predicted
      if ( state == 0 || state == 1) {
        miss++
      }
      //if state is not ST, increment
      if (state != 3) {
        state++
      }

      GBH |= int (math.Pow(2, float64(N)))
    } else { //Branch not taken
      //Check if it was Miss Predicted
      if ( state == 2 || state == 3) {
        miss++
      }
      //if state is SN, decrement
      if (state != 0) {
        state--
      }
    }
    GBH >>= 1 //remove LSB
    a[i] = state //update state of array a
    branch_num++ //increase the branch
  }
  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  missR := (float32(miss)/float32(branch_num)) * 100 //gives percentage for miss ratio

  fmt.Printf("%d %d %.2f", M, N, missR)
  fmt.Println()
}
