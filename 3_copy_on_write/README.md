Реализовать COW (Copy-On-Write) буфер.

// Предположим, что все будут производить копирование буфера только с использованием метода Clone()
```go
type COWBuffer struct { ... }

func NewCOWBuffer(data []byte)                         // создать буффер с определенными данными
func (b *COWBuffer) Clone() COWBuffer                  // создать новую копию буфера
func (b *COWBuffer) Close()                            // перестать использовать копию буффера
func (b *COWBuffer) Update(index int, value byte) bool // изменить определенный байт в буффере
func (b *COWBuffer) String() string                    // сконвертировать буффер в строку
``` 