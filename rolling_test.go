package main

import (
	"crypto/sha256"
	"reflect"
	"testing"
)

func TestRollingHashDiffSameContents(t *testing.T) {
	original := []byte("Testx")
	updated := []byte("Testx")
	chunkSize := 1
	hashFunc := sha256.New

	expected := fileDelta{
		Reusables: []fileChunk{
			{
				Index: 0,
				Data:  []byte("T"),
			},
			{
				Index: 1,
				Data:  []byte("e"),
			},
			{
				Index: 2,
				Data:  []byte("s"),
			},
			{
				Index: 3,
				Data:  []byte("t"),
			},
			{
				Index: 4,
				Data:  []byte("x"),
			},
		},
		Modified: nil,
		Removed:  nil,
	}
	result := rollingHashDiff(original, updated, hashFunc, chunkSize)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("For original %v and updated %v, expected %v but got %v", original, updated, expected, result)
	}
}

func TestRollingHashDiffUpdatedContent(t *testing.T) {
	original := []byte("Testx")
	updated := []byte("Testz")
	chunkSize := 1
	hashFunc := sha256.New
	expected := fileDelta{
		Reusables: []fileChunk{
			{
				Index: 0,
				Data:  []byte("T"),
			},
			{
				Index: 1,
				Data:  []byte("e"),
			},
			{
				Index: 2,
				Data:  []byte("s"),
			},
			{
				Index: 3,
				Data:  []byte("t"),
			},
		},
		Modified: []fileChunk{
			{
				Index: 4,
				Data:  []byte("z"),
			},
		},
		Removed: []fileChunk{
			{
				Index: 4,
				Data:  []byte("x"),
			},
		},
	}
	result := rollingHashDiff(original, updated, hashFunc, chunkSize)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("For original %v and updated %v, expected %v but got %v", original, updated, expected, result)
	}
}

func TestRollingHashDiffUpdatedContent2(t *testing.T) {
	original := []byte("Testx")
	updated := []byte("Tessz")
	chunkSize := 1
	hashFunc := sha256.New
	expected := fileDelta{
		Reusables: []fileChunk{
			{
				Index: 0,
				Data:  []byte("T"),
			},
			{
				Index: 1,
				Data:  []byte("e"),
			},
			{
				Index: 2,
				Data:  []byte("s"),
			},
		},
		Modified: []fileChunk{
			{
				Index: 3,
				Data:  []byte("s"),
			},
			{
				Index: 4,
				Data:  []byte("z"),
			},
		},
		Removed: []fileChunk{
			{
				Index: 3,
				Data:  []byte("t"),
			},
			{
				Index: 4,
				Data:  []byte("x"),
			},
		},
	}
	result := rollingHashDiff(original, updated, hashFunc, chunkSize)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("For original %v and updated %v, expected %v but got %v", original, updated, expected, result)
	}
}

func TestRollingHashChunkSize(t *testing.T) {
	original := []byte("Same files")
	updated := []byte("Emas files")
	chunkSize := 5
	hashFunc := sha256.New
	expected := fileDelta{
		Reusables: []fileChunk{
			{
				Index: 1,
				Data:  []byte("files"),
			},
		},
		Modified: []fileChunk{
			{
				Index: 0,
				Data:  []byte("Emas "),
			},
		},
		Removed: []fileChunk{
			{
				Index: 0,
				Data:  []byte("Same "),
			},
		},
	}
	result := rollingHashDiff(original, updated, hashFunc, chunkSize)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("For original %v and updated %v, expected %v but got %v", original, updated, expected, result)
	}
}

func TestRollingHashRemovedChunk(t *testing.T) {
	original := []byte("Same sure files")
	updated := []byte("Same files")
	chunkSize := 5
	hashFunc := sha256.New
	expected := fileDelta{
		Reusables: []fileChunk{
			{
				Index: 0,
				Data:  []byte("Same "),
			},
		},
		Modified: []fileChunk{
			{
				Index: 1,
				Data:  []byte("files"),
			},
			{
				Index: 2,
				Data:  []byte(""),
			},
		},
		Removed: []fileChunk{
			{
				Index: 1,
				Data:  []byte("sure "),
			},
			{
				Index: 2,
				Data:  []byte("files"),
			},
		},
	}
	result := rollingHashDiff(original, updated, hashFunc, chunkSize)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("For original %v and updated %v, expected %v but got %v", original, updated, expected, result)
	}
}

func TestRollingHashAddedChunk(t *testing.T) {
	original := []byte("Same files")
	updated := []byte("Same files, larger!!")
	chunkSize := 10
	hashFunc := sha256.New
	expected := fileDelta{
		Reusables: []fileChunk{
			{
				Index: 0,
				Data:  []byte("Same files"),
			},
		},
		Modified: []fileChunk{
			{
				Index: 1,
				Data:  []byte(", larger!!"),
			},
		},
		Removed: []fileChunk{
			{
				Index: 1,
				Data:  []byte(""),
			},
		},
	}
	result := rollingHashDiff(original, updated, hashFunc, chunkSize)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("For original %v and updated %v, expected %v but got %v", original, updated, expected, result)
	}
}
