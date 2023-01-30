package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"
)

// fileDelta represents the output of a diffing operation
type fileDelta struct {
	Reusables []fileChunk
	Modified  []fileChunk
	Removed   []fileChunk
}

func (fc fileDelta) String() string {
	return fmt.Sprintf("Reusables: %v, Modified: %v, Removed: %v", fc.Reusables, fc.Modified, fc.Removed)
}

// fileChunk represents a given chunk on a file 
type fileChunk struct {
	Index int
	Data  []byte
}

func (fc fileChunk) String() string {
	return fmt.Sprintf("Index: %v, Data: %s", fc.Index, string(fc.Data))
}

// rollingHashDiff compares the original and updated version of a file using a rolling hash and returns a description of 
// the chunks that can be reused and the chunks that have been added, modified, or removed
func rollingHashDiff(originalBytes, updatedBytes []byte, hashFunc func() hash.Hash, chunkSize int) fileDelta {
	//originalHash := hashFunc()
	//updatedHash := hashFunc()
	var delta []fileChunk
	var reusables []fileChunk
	var removed []fileChunk

	originalChunks := chunkData(originalBytes, chunkSize)
	updatedChunks := chunkData(updatedBytes, chunkSize)

	for i := len(originalChunks); i < len(updatedChunks); i++ {
		originalChunks = append(originalChunks, []byte{})
	}

	for i, originalChunk := range originalChunks {
		updatedChunk := updatedChunks[i]
		if bytes.Equal(originalChunk, updatedChunk) {
			reusables = append(reusables, fileChunk{Index: i, Data: originalChunk})
		} else {
			delta = append(delta, fileChunk{Index: i, Data: updatedChunk})
			removed = append(removed, fileChunk{Index: i, Data: originalChunk})
		}
	}

	return fileDelta{Reusables: reusables, Modified: delta, Removed: removed}
}

// chunkData splits the original data into chunks according to the chunkSize parameter
func chunkData(data []byte, chunkSize int) [][]byte {
	var chunks [][]byte
	
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
}

func main() {
	result1 := rollingHashDiff([]byte("Testx"), []byte("Testx"), sha256.New, 1)
	fmt.Println(result1)
	result2 := rollingHashDiff([]byte("Testx"), []byte("Testz"), sha256.New, 1)
	fmt.Println(result2)
	result3 := rollingHashDiff([]byte("Testx"), []byte("Tessz"), sha256.New, 1)
	fmt.Println(result3)
}
