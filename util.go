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
