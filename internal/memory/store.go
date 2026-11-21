package memory

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"syscall"
)

type NodeRecord struct {
	UUID     string
	Type     string // MODULE, ENTITY, LAW, FLOW
	Name     string
	Parent   string
	Children []string
	Props    map[string]string // Propiedades e inyecciones
}

type MemoryStore struct {
	filePath string
	Nodes    map[string]*NodeRecord
}

func init() {
	gob.Register(&NodeRecord{})
}

func NewMemoryStore(filePath string) *MemoryStore {
	return &MemoryStore{
		filePath: filePath,
		Nodes:    make(map[string]*NodeRecord),
	}
}

// Persiste la estructura en disco (.axm) con encoding/gob
func (ms *MemoryStore) Save() error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(ms.Nodes); err != nil {
		return err
	}

	return os.WriteFile(ms.filePath, buf.Bytes(), 0644)
}

// Carga directa mediante syscall.Mmap para cero latencia de lectura
func (ms *MemoryStore) LoadMmap() error {
	file, err := os.OpenFile(ms.filePath, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil || info.Size() == 0 {
		return nil
	}

	// Conversión segura de int64 a int para syscall
	size := int(info.Size())

	// Memory-Mapped File a nivel de Kernel (Debian/Posix)
	mmapData, err := syscall.Mmap(int(file.Fd()), 0, size, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return fmt.Errorf("error mmap: %v", err)
	}
	defer syscall.Munmap(mmapData)

	buf := bytes.NewBuffer(mmapData)
	decoder := gob.NewDecoder(buf)

	var loadedNodes map[string]*NodeRecord
	if err := decoder.Decode(&loadedNodes); err != nil {
		return err
	}

	ms.Nodes = loadedNodes
	return nil
}
