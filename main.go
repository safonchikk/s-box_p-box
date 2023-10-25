package main

import (
	"fmt"
)

func getSTable(key int) []byte {
	table := [][]byte{
		{0b0010, 0b1100, 0b0100, 0b0001, 0b0111, 0b1010, 0b1011, 0b0110, 0b1000, 0b0101, 0b0011, 0b1111, 0b1101, 0b0000, 0b1110, 0b1001},
		{0b1110, 0b1011, 0b0010, 0b1100, 0b0100, 0b0111, 0b1101, 0b0001, 0b0101, 0b0000, 0b1111, 0b1010, 0b0011, 0b1001, 0b1000, 0b0110},
		{0b0100, 0b0010, 0b0001, 0b1011, 0b1010, 0b1101, 0b0111, 0b1000, 0b1111, 0b1001, 0b1100, 0b0101, 0b0110, 0b0011, 0b0000, 0b1110},
		{0b1011, 0b1000, 0b1100, 0b0111, 0b0001, 0b1110, 0b0010, 0b1101, 0b0110, 0b1111, 0b0000, 0b1001, 0b1010, 0b0100, 0b0101, 0b0011},
	}
	return table[key]
}

func SBoxEncrypt(input, output *[]byte, key int) { //encrypts byte array, key chooses one of rows in table above
	sTable := getSTable(key)
	firstHalf := byte(0)
	secondHalf := byte(0)
	for i, block := range *input {
		firstHalf = block >> 4
		secondHalf = block & 15
		firstHalf = sTable[firstHalf]
		secondHalf = sTable[secondHalf]
		(*output)[i] = (firstHalf << 4) + secondHalf
	}
}

func find(ar []byte, s byte) byte {
	for i, el := range ar {
		if el == s {
			return byte(i)
		}
	}
	return 0
}

func SBoxDecrypt(input, output *[]byte, key int) { //decrypts byte array back, needs same key as encoder
	sTable := getSTable(key) //we can have either reverse table or search in original one
	firstHalf := byte(0)     //reverse table is faster but needs separate function to calculate it
	secondHalf := byte(0)
	for i, block := range *input {
		firstHalf = block >> 4
		secondHalf = block & 15
		firstHalf = find(sTable, block>>4)
		secondHalf = find(sTable, block&15)
		(*output)[i] = (firstHalf << 4) + secondHalf
	}
}

func getBits(dest *[]byte, input byte) { //turns byte into array of bits, small first
	for i := 0; i < 8; i++ {
		(*dest)[i] = (input >> i) & 1
	}
}

func numFromBits(input *[]byte) byte { //turns array of bits into number
	res := byte(0)
	for i := 0; i < 8; i++ {
		res += (*input)[i] << i
	}
	return res
}

func PBoxEncrypt(input, output *[]byte) { //encrypts array of bytes
	pTable := []byte{3, 7, 0, 1, 6, 2, 5, 4}
	bits := make([]byte, 8)
	outBits := make([]byte, 8)
	for i, el := range *input {
		getBits(&bits, el)
		for j := 0; j < 8; j++ {
			outBits[j] = bits[pTable[j]]
		}
		(*output)[i] = numFromBits(&outBits)
	}
}

func PBoxDecrypt(input, output *[]byte) { //decrypts array of bytes
	pTable := []byte{2, 3, 5, 0, 7, 6, 4, 1}
	bits := make([]byte, 8)
	outBits := make([]byte, 8)
	for i, el := range *input {
		getBits(&bits, el)
		for j := 0; j < 8; j++ {
			outBits[j] = bits[pTable[j]]
		}
		(*output)[i] = numFromBits(&outBits)
	}
}

func main() {
	input := []byte{
		195, 16, 74, 83, 55, 4, 61, 100,
	}
	sOutput := make([]byte, len(input))
	sOutput2 := make([]byte, len(input))

	SBoxEncrypt(&input, &sOutput, 2)
	SBoxDecrypt(&sOutput, &sOutput2, 2)

	fmt.Print("input:              ")
	fmt.Println(input)
	fmt.Print("encrypted by s-boxes: ")
	fmt.Println(sOutput)
	fmt.Print("decrypted by s-boxes: ")
	fmt.Println(sOutput2)

	pOutput := make([]byte, len(input))
	pOutput2 := make([]byte, len(input))

	PBoxEncrypt(&input, &pOutput)
	PBoxDecrypt(&pOutput, &pOutput2)

	fmt.Print("encrypted by p-boxes: ")
	fmt.Println(pOutput)
	fmt.Print("decrypted by p-boxes: ")
	fmt.Println(pOutput2)
}
