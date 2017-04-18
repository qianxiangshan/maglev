//100个服务以内的一致性hash算法
package maglev

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	mTableSize = 10273
)

type metadata struct {
	//backendserver size
	n                 int
	lookuptable       [mTableSize]int
	indexservicetable map[uint16]string
	servicenames      []string
}

type Maglev struct {
	metadata
}

//input servicename
func (maghash *Maglev) Init(servcenames []string) error {

	maghash.n = len(servcenames)
	if maghash.n == 0 || maghash.n > 100 {
		return errors.New("servicename countn is 0 or bigger than 100")
	}

	maghash.servicenames = servcenames

	maghash.indexservicetable = make(map[uint16]string)

	for index, name := range servcenames {
		maghash.indexservicetable[uint16(index)] = name
	}

	var permutition [][]int
	permutition = make([][]int, mTableSize)
	for i := 0; i < mTableSize; i++ {
		permutition[i] = make([]int, maghash.n)
	}

	maghash.genPermutation(permutition)
	//fmt.Printf("permutation %v\n", permutition)
	maghash.genLookupTable(permutition)
	//fmt.Printf("lookuptable %v\n", maghash.lookuptable)

	return nil
}

//获取key对应的servicename
func (maghash *Maglev) Get(key string) (servicename string) {

	index := genOffset(key) % mTableSize

	serviceindex := maghash.lookuptable[index]

	return maghash.indexservicetable[uint16(serviceindex)]
}

func (maghash *Maglev) genLookupTable(permutition [][]int) {
	var i, j int
	var next = make([]int, maghash.n)

	for ; j < mTableSize; j++ {
		maghash.lookuptable[j] = -1
	}

	var n int = 0
	var c int
	i, j = 0, 0
	for {
		for i = 0; i < maghash.n; i++ {
			c = permutition[next[i]][i]
			for maghash.lookuptable[c] >= 0 {
				next[i] = next[i] + 1
				c = permutition[next[i]][i]
			}
			maghash.lookuptable[c] = i
			next[i] = next[i] + 1
			n = n + 1
			if n == mTableSize {
				fmt.Println("gen lookuptable over")
				return
			}
		}
	}

}

//生成permutation表
func (maghash *Maglev) genPermutation(permutition [][]int) {

	for i := 0; i < maghash.n; i++ {
		for j := 0; j < mTableSize; j++ {
			permutition[j][i] = (genOffset(maghash.servicenames[i]) + j*genSkip(maghash.servicenames[i])) % mTableSize
		}
	}
}

func genSkip(servicename string) int {
	sum := sha256.Sum256([]byte(servicename))

	reader := bytes.NewReader(sum[:8])

	var skip uint64
	err := binary.Read(reader, binary.BigEndian, &skip)
	if err != nil {
		panic("binary to uint64 error" + err.Error())
	}

	skip = skip % mTableSize

	return int(skip)

}
func genOffset(servicename string) int {
	sum := md5.Sum([]byte(servicename))
	reader := bytes.NewReader(sum[:8])

	var offset uint64
	err := binary.Read(reader, binary.BigEndian, &offset)
	if err != nil {
		panic("binary to uint64 error" + err.Error())
	}

	offset = offset%(mTableSize-1) + 1

	return int(offset)

}
