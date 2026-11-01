// Package catalog provides large module
package catalog

import (
	"fmt"
	"sync"
)

// Entity1
type Entity1 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity1(id,name string) *Entity1 { return &Entity1{ID:id,Name:name} }
func (e *Entity1) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity1) String() string { return fmt.Sprintf("Entity1 %s %s", e.ID, e.Name) }

// Entity2
type Entity2 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity2(id,name string) *Entity2 { return &Entity2{ID:id,Name:name} }
func (e *Entity2) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity2) String() string { return fmt.Sprintf("Entity2 %s %s", e.ID, e.Name) }

// Entity3
type Entity3 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity3(id,name string) *Entity3 { return &Entity3{ID:id,Name:name} }
func (e *Entity3) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity3) String() string { return fmt.Sprintf("Entity3 %s %s", e.ID, e.Name) }

// Entity4
type Entity4 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity4(id,name string) *Entity4 { return &Entity4{ID:id,Name:name} }
func (e *Entity4) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity4) String() string { return fmt.Sprintf("Entity4 %s %s", e.ID, e.Name) }

// Entity5
type Entity5 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity5(id,name string) *Entity5 { return &Entity5{ID:id,Name:name} }
func (e *Entity5) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity5) String() string { return fmt.Sprintf("Entity5 %s %s", e.ID, e.Name) }

// Entity6
type Entity6 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity6(id,name string) *Entity6 { return &Entity6{ID:id,Name:name} }
func (e *Entity6) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity6) String() string { return fmt.Sprintf("Entity6 %s %s", e.ID, e.Name) }

// Entity7
type Entity7 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity7(id,name string) *Entity7 { return &Entity7{ID:id,Name:name} }
func (e *Entity7) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity7) String() string { return fmt.Sprintf("Entity7 %s %s", e.ID, e.Name) }

// Entity8
type Entity8 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity8(id,name string) *Entity8 { return &Entity8{ID:id,Name:name} }
func (e *Entity8) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity8) String() string { return fmt.Sprintf("Entity8 %s %s", e.ID, e.Name) }

// Entity9
type Entity9 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity9(id,name string) *Entity9 { return &Entity9{ID:id,Name:name} }
func (e *Entity9) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity9) String() string { return fmt.Sprintf("Entity9 %s %s", e.ID, e.Name) }

// Entity10
type Entity10 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity10(id,name string) *Entity10 { return &Entity10{ID:id,Name:name} }
func (e *Entity10) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity10) String() string { return fmt.Sprintf("Entity10 %s %s", e.ID, e.Name) }

// Entity11
type Entity11 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity11(id,name string) *Entity11 { return &Entity11{ID:id,Name:name} }
func (e *Entity11) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity11) String() string { return fmt.Sprintf("Entity11 %s %s", e.ID, e.Name) }

// Entity12
type Entity12 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity12(id,name string) *Entity12 { return &Entity12{ID:id,Name:name} }
func (e *Entity12) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity12) String() string { return fmt.Sprintf("Entity12 %s %s", e.ID, e.Name) }

// Entity13
type Entity13 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity13(id,name string) *Entity13 { return &Entity13{ID:id,Name:name} }
func (e *Entity13) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity13) String() string { return fmt.Sprintf("Entity13 %s %s", e.ID, e.Name) }

// Entity14
type Entity14 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity14(id,name string) *Entity14 { return &Entity14{ID:id,Name:name} }
func (e *Entity14) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity14) String() string { return fmt.Sprintf("Entity14 %s %s", e.ID, e.Name) }

// Entity15
type Entity15 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity15(id,name string) *Entity15 { return &Entity15{ID:id,Name:name} }
func (e *Entity15) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity15) String() string { return fmt.Sprintf("Entity15 %s %s", e.ID, e.Name) }

// Entity16
type Entity16 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity16(id,name string) *Entity16 { return &Entity16{ID:id,Name:name} }
func (e *Entity16) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity16) String() string { return fmt.Sprintf("Entity16 %s %s", e.ID, e.Name) }

// Entity17
type Entity17 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity17(id,name string) *Entity17 { return &Entity17{ID:id,Name:name} }
func (e *Entity17) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity17) String() string { return fmt.Sprintf("Entity17 %s %s", e.ID, e.Name) }

// Entity18
type Entity18 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity18(id,name string) *Entity18 { return &Entity18{ID:id,Name:name} }
func (e *Entity18) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity18) String() string { return fmt.Sprintf("Entity18 %s %s", e.ID, e.Name) }

// Entity19
type Entity19 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity19(id,name string) *Entity19 { return &Entity19{ID:id,Name:name} }
func (e *Entity19) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity19) String() string { return fmt.Sprintf("Entity19 %s %s", e.ID, e.Name) }

// Entity20
type Entity20 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity20(id,name string) *Entity20 { return &Entity20{ID:id,Name:name} }
func (e *Entity20) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity20) String() string { return fmt.Sprintf("Entity20 %s %s", e.ID, e.Name) }

// Entity21
type Entity21 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity21(id,name string) *Entity21 { return &Entity21{ID:id,Name:name} }
func (e *Entity21) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity21) String() string { return fmt.Sprintf("Entity21 %s %s", e.ID, e.Name) }

// Entity22
type Entity22 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity22(id,name string) *Entity22 { return &Entity22{ID:id,Name:name} }
func (e *Entity22) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity22) String() string { return fmt.Sprintf("Entity22 %s %s", e.ID, e.Name) }

// Entity23
type Entity23 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity23(id,name string) *Entity23 { return &Entity23{ID:id,Name:name} }
func (e *Entity23) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity23) String() string { return fmt.Sprintf("Entity23 %s %s", e.ID, e.Name) }

// Entity24
type Entity24 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity24(id,name string) *Entity24 { return &Entity24{ID:id,Name:name} }
func (e *Entity24) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity24) String() string { return fmt.Sprintf("Entity24 %s %s", e.ID, e.Name) }

// Entity25
type Entity25 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity25(id,name string) *Entity25 { return &Entity25{ID:id,Name:name} }
func (e *Entity25) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity25) String() string { return fmt.Sprintf("Entity25 %s %s", e.ID, e.Name) }

// Entity26
type Entity26 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity26(id,name string) *Entity26 { return &Entity26{ID:id,Name:name} }
func (e *Entity26) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity26) String() string { return fmt.Sprintf("Entity26 %s %s", e.ID, e.Name) }

// Entity27
type Entity27 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity27(id,name string) *Entity27 { return &Entity27{ID:id,Name:name} }
func (e *Entity27) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity27) String() string { return fmt.Sprintf("Entity27 %s %s", e.ID, e.Name) }

// Entity28
type Entity28 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity28(id,name string) *Entity28 { return &Entity28{ID:id,Name:name} }
func (e *Entity28) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity28) String() string { return fmt.Sprintf("Entity28 %s %s", e.ID, e.Name) }

// Entity29
type Entity29 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity29(id,name string) *Entity29 { return &Entity29{ID:id,Name:name} }
func (e *Entity29) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity29) String() string { return fmt.Sprintf("Entity29 %s %s", e.ID, e.Name) }

// Entity30
type Entity30 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity30(id,name string) *Entity30 { return &Entity30{ID:id,Name:name} }
func (e *Entity30) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity30) String() string { return fmt.Sprintf("Entity30 %s %s", e.ID, e.Name) }

// Entity31
type Entity31 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity31(id,name string) *Entity31 { return &Entity31{ID:id,Name:name} }
func (e *Entity31) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity31) String() string { return fmt.Sprintf("Entity31 %s %s", e.ID, e.Name) }

// Entity32
type Entity32 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity32(id,name string) *Entity32 { return &Entity32{ID:id,Name:name} }
func (e *Entity32) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity32) String() string { return fmt.Sprintf("Entity32 %s %s", e.ID, e.Name) }

// Entity33
type Entity33 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity33(id,name string) *Entity33 { return &Entity33{ID:id,Name:name} }
func (e *Entity33) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity33) String() string { return fmt.Sprintf("Entity33 %s %s", e.ID, e.Name) }

// Entity34
type Entity34 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity34(id,name string) *Entity34 { return &Entity34{ID:id,Name:name} }
func (e *Entity34) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity34) String() string { return fmt.Sprintf("Entity34 %s %s", e.ID, e.Name) }

// Entity35
type Entity35 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity35(id,name string) *Entity35 { return &Entity35{ID:id,Name:name} }
func (e *Entity35) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity35) String() string { return fmt.Sprintf("Entity35 %s %s", e.ID, e.Name) }

// Entity36
type Entity36 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity36(id,name string) *Entity36 { return &Entity36{ID:id,Name:name} }
func (e *Entity36) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity36) String() string { return fmt.Sprintf("Entity36 %s %s", e.ID, e.Name) }

// Entity37
type Entity37 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity37(id,name string) *Entity37 { return &Entity37{ID:id,Name:name} }
func (e *Entity37) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity37) String() string { return fmt.Sprintf("Entity37 %s %s", e.ID, e.Name) }

// Entity38
type Entity38 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity38(id,name string) *Entity38 { return &Entity38{ID:id,Name:name} }
func (e *Entity38) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity38) String() string { return fmt.Sprintf("Entity38 %s %s", e.ID, e.Name) }

// Entity39
type Entity39 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity39(id,name string) *Entity39 { return &Entity39{ID:id,Name:name} }
func (e *Entity39) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity39) String() string { return fmt.Sprintf("Entity39 %s %s", e.ID, e.Name) }

// Entity40
type Entity40 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewEntity40(id,name string) *Entity40 { return &Entity40{ID:id,Name:name} }
func (e *Entity40) Rename(name string) { e.mu.Lock(); defer e.mu.Unlock(); e.Name=name }
func (e *Entity40) String() string { return fmt.Sprintf("Entity40 %s %s", e.ID, e.Name) }


