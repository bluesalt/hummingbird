package objectserver

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteReadMetadata(t *testing.T) {

	data := map[string]interface{}{
		strings.Repeat("la", 5):    strings.Repeat("la", 30),
		strings.Repeat("moo", 500): strings.Repeat("moo", 300),
	}
	testFile, err := ioutil.TempFile("/tmp", "backend_test")
	defer testFile.Close()
	defer os.Remove(testFile.Name())
	assert.Equal(t, err, nil)
	WriteMetadata(testFile.Fd(), data)
	checkData := map[interface{}]interface{}{
		strings.Repeat("la", 5):    strings.Repeat("la", 30),
		strings.Repeat("moo", 500): strings.Repeat("moo", 300),
	}
	readData, err := ReadMetadata(testFile.Name())
	assert.Equal(t, err, nil)
	assert.True(t, reflect.DeepEqual(checkData, readData))

	readData, err = ReadMetadata(testFile.Fd())
	assert.Equal(t, err, nil)
	assert.True(t, reflect.DeepEqual(checkData, readData))
}

func TestGetHashes(t *testing.T) {
	driveRoot, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(driveRoot)
	os.MkdirAll(driveRoot+"/sda/objects/1/abc/fffffffffffffffffffffffffffffabc", 0777)
	os.MkdirAll(driveRoot+"/sda/objects/1/abc/00000000000000000000000000000abc", 0777)
	f, _ := os.Create(driveRoot + "/sda/objects/1/abc/fffffffffffffffffffffffffffffabc/12345.data")
	defer f.Close()
	f, _ = os.Create(driveRoot + "/sda/objects/1/abc/00000000000000000000000000000abc/67890.data")
	defer f.Close()

	hashes, err := GetHashes(driveRoot, "sda", "1", nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, "b1589029b7db9d01347caece2159d588", hashes["abc"])

	// write a new file there
	f, _ = os.Create(driveRoot + "/sda/objects/1/abc/00000000000000000000000000000abc/99999.meta")
	f.Close()

	// make sure hash for "abc" isn't recalculated yet.
	hashes, err = GetHashes(driveRoot, "sda", "1", nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, "b1589029b7db9d01347caece2159d588", hashes["abc"])

	// force recalculate of "abc"
	hashes, err = GetHashes(driveRoot, "sda", "1", []string{"abc"}, nil)
	assert.Nil(t, err)
	assert.Equal(t, "8834e84467693c2e8f670f4afbea5334", hashes["abc"])
}