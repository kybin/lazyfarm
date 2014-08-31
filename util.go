package main

func removeDuplicates(a []int) []int {
        result := []int{}
        seen := map[int]int{}
        for _, val := range a {
                if _, ok := seen[val]; !ok {
                        result = append(result, val)
                        seen[val] = val
                }
        }
        return result
}

type intSlice []int

func (slice intSlice) pos(value int) int {
    for p, v := range slice {
        if (v == value) {
            return p
        }
    }
    return -1
}
