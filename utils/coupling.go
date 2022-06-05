package utils
import (
    "math/rand"
	"regexp"
)

// カップリングする
func Coupling(lines [][]string, coupling [][]string) [][]string {

	var couple []string
	last := false
	
    for i := 0; i < len(lines); i++{
        if len(lines[i]) <= 1 { last = true }
		line := lines[i]
		c := rand.Intn(len(line))
		couple = append(couple, line[c])
		
		line[c] = line[len(line)-1] 
		lines[i] = line[:len(line)-1]
    }
	coupling = append(coupling, couple)
	if !last { coupling = Coupling(lines, coupling) }

	return coupling
}
// 与えられた列をカンマ区切りで配列に格納して２次元配列にする
func GetItemSets(data []string, splitter string) [][]string {
    var lines [][]string
    for i := 0; i < len(data); i++{
        line := data[i]
		t := regexp.MustCompile(splitter).Split(line, -1)
		lines = append(lines, t)
    }
	return lines
}
