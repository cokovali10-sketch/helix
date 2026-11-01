// Package catalog provides large module 5
package catalog

import (
	"fmt"
	"sync"
)

// CatalogEntry1
type CatalogEntry1 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry1(id,name string) *CatalogEntry1 { return &CatalogEntry1{ID:id,Name:name} }
func (c *CatalogEntry1) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry1) String() string { return fmt.Sprintf("CatalogEntry1 %s %s", c.ID, c.Name) }

// CatalogEntry2
type CatalogEntry2 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry2(id,name string) *CatalogEntry2 { return &CatalogEntry2{ID:id,Name:name} }
func (c *CatalogEntry2) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry2) String() string { return fmt.Sprintf("CatalogEntry2 %s %s", c.ID, c.Name) }

// CatalogEntry3
type CatalogEntry3 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry3(id,name string) *CatalogEntry3 { return &CatalogEntry3{ID:id,Name:name} }
func (c *CatalogEntry3) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry3) String() string { return fmt.Sprintf("CatalogEntry3 %s %s", c.ID, c.Name) }

// CatalogEntry4
type CatalogEntry4 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry4(id,name string) *CatalogEntry4 { return &CatalogEntry4{ID:id,Name:name} }
func (c *CatalogEntry4) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry4) String() string { return fmt.Sprintf("CatalogEntry4 %s %s", c.ID, c.Name) }

// CatalogEntry5
type CatalogEntry5 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry5(id,name string) *CatalogEntry5 { return &CatalogEntry5{ID:id,Name:name} }
func (c *CatalogEntry5) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry5) String() string { return fmt.Sprintf("CatalogEntry5 %s %s", c.ID, c.Name) }

// CatalogEntry6
type CatalogEntry6 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry6(id,name string) *CatalogEntry6 { return &CatalogEntry6{ID:id,Name:name} }
func (c *CatalogEntry6) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry6) String() string { return fmt.Sprintf("CatalogEntry6 %s %s", c.ID, c.Name) }

// CatalogEntry7
type CatalogEntry7 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry7(id,name string) *CatalogEntry7 { return &CatalogEntry7{ID:id,Name:name} }
func (c *CatalogEntry7) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry7) String() string { return fmt.Sprintf("CatalogEntry7 %s %s", c.ID, c.Name) }

// CatalogEntry8
type CatalogEntry8 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry8(id,name string) *CatalogEntry8 { return &CatalogEntry8{ID:id,Name:name} }
func (c *CatalogEntry8) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry8) String() string { return fmt.Sprintf("CatalogEntry8 %s %s", c.ID, c.Name) }

// CatalogEntry9
type CatalogEntry9 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry9(id,name string) *CatalogEntry9 { return &CatalogEntry9{ID:id,Name:name} }
func (c *CatalogEntry9) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry9) String() string { return fmt.Sprintf("CatalogEntry9 %s %s", c.ID, c.Name) }

// CatalogEntry10
type CatalogEntry10 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry10(id,name string) *CatalogEntry10 { return &CatalogEntry10{ID:id,Name:name} }
func (c *CatalogEntry10) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry10) String() string { return fmt.Sprintf("CatalogEntry10 %s %s", c.ID, c.Name) }

// CatalogEntry11
type CatalogEntry11 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry11(id,name string) *CatalogEntry11 { return &CatalogEntry11{ID:id,Name:name} }
func (c *CatalogEntry11) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry11) String() string { return fmt.Sprintf("CatalogEntry11 %s %s", c.ID, c.Name) }

// CatalogEntry12
type CatalogEntry12 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry12(id,name string) *CatalogEntry12 { return &CatalogEntry12{ID:id,Name:name} }
func (c *CatalogEntry12) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry12) String() string { return fmt.Sprintf("CatalogEntry12 %s %s", c.ID, c.Name) }

// CatalogEntry13
type CatalogEntry13 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry13(id,name string) *CatalogEntry13 { return &CatalogEntry13{ID:id,Name:name} }
func (c *CatalogEntry13) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry13) String() string { return fmt.Sprintf("CatalogEntry13 %s %s", c.ID, c.Name) }

// CatalogEntry14
type CatalogEntry14 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry14(id,name string) *CatalogEntry14 { return &CatalogEntry14{ID:id,Name:name} }
func (c *CatalogEntry14) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry14) String() string { return fmt.Sprintf("CatalogEntry14 %s %s", c.ID, c.Name) }

// CatalogEntry15
type CatalogEntry15 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry15(id,name string) *CatalogEntry15 { return &CatalogEntry15{ID:id,Name:name} }
func (c *CatalogEntry15) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry15) String() string { return fmt.Sprintf("CatalogEntry15 %s %s", c.ID, c.Name) }

// CatalogEntry16
type CatalogEntry16 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry16(id,name string) *CatalogEntry16 { return &CatalogEntry16{ID:id,Name:name} }
func (c *CatalogEntry16) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry16) String() string { return fmt.Sprintf("CatalogEntry16 %s %s", c.ID, c.Name) }

// CatalogEntry17
type CatalogEntry17 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry17(id,name string) *CatalogEntry17 { return &CatalogEntry17{ID:id,Name:name} }
func (c *CatalogEntry17) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry17) String() string { return fmt.Sprintf("CatalogEntry17 %s %s", c.ID, c.Name) }

// CatalogEntry18
type CatalogEntry18 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry18(id,name string) *CatalogEntry18 { return &CatalogEntry18{ID:id,Name:name} }
func (c *CatalogEntry18) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry18) String() string { return fmt.Sprintf("CatalogEntry18 %s %s", c.ID, c.Name) }

// CatalogEntry19
type CatalogEntry19 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry19(id,name string) *CatalogEntry19 { return &CatalogEntry19{ID:id,Name:name} }
func (c *CatalogEntry19) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry19) String() string { return fmt.Sprintf("CatalogEntry19 %s %s", c.ID, c.Name) }

// CatalogEntry20
type CatalogEntry20 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry20(id,name string) *CatalogEntry20 { return &CatalogEntry20{ID:id,Name:name} }
func (c *CatalogEntry20) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry20) String() string { return fmt.Sprintf("CatalogEntry20 %s %s", c.ID, c.Name) }

// CatalogEntry21
type CatalogEntry21 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry21(id,name string) *CatalogEntry21 { return &CatalogEntry21{ID:id,Name:name} }
func (c *CatalogEntry21) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry21) String() string { return fmt.Sprintf("CatalogEntry21 %s %s", c.ID, c.Name) }

// CatalogEntry22
type CatalogEntry22 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry22(id,name string) *CatalogEntry22 { return &CatalogEntry22{ID:id,Name:name} }
func (c *CatalogEntry22) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry22) String() string { return fmt.Sprintf("CatalogEntry22 %s %s", c.ID, c.Name) }

// CatalogEntry23
type CatalogEntry23 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry23(id,name string) *CatalogEntry23 { return &CatalogEntry23{ID:id,Name:name} }
func (c *CatalogEntry23) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry23) String() string { return fmt.Sprintf("CatalogEntry23 %s %s", c.ID, c.Name) }

// CatalogEntry24
type CatalogEntry24 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry24(id,name string) *CatalogEntry24 { return &CatalogEntry24{ID:id,Name:name} }
func (c *CatalogEntry24) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry24) String() string { return fmt.Sprintf("CatalogEntry24 %s %s", c.ID, c.Name) }

// CatalogEntry25
type CatalogEntry25 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry25(id,name string) *CatalogEntry25 { return &CatalogEntry25{ID:id,Name:name} }
func (c *CatalogEntry25) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry25) String() string { return fmt.Sprintf("CatalogEntry25 %s %s", c.ID, c.Name) }

// CatalogEntry26
type CatalogEntry26 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry26(id,name string) *CatalogEntry26 { return &CatalogEntry26{ID:id,Name:name} }
func (c *CatalogEntry26) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry26) String() string { return fmt.Sprintf("CatalogEntry26 %s %s", c.ID, c.Name) }

// CatalogEntry27
type CatalogEntry27 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry27(id,name string) *CatalogEntry27 { return &CatalogEntry27{ID:id,Name:name} }
func (c *CatalogEntry27) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry27) String() string { return fmt.Sprintf("CatalogEntry27 %s %s", c.ID, c.Name) }

// CatalogEntry28
type CatalogEntry28 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry28(id,name string) *CatalogEntry28 { return &CatalogEntry28{ID:id,Name:name} }
func (c *CatalogEntry28) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry28) String() string { return fmt.Sprintf("CatalogEntry28 %s %s", c.ID, c.Name) }

// CatalogEntry29
type CatalogEntry29 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry29(id,name string) *CatalogEntry29 { return &CatalogEntry29{ID:id,Name:name} }
func (c *CatalogEntry29) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry29) String() string { return fmt.Sprintf("CatalogEntry29 %s %s", c.ID, c.Name) }

// CatalogEntry30
type CatalogEntry30 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry30(id,name string) *CatalogEntry30 { return &CatalogEntry30{ID:id,Name:name} }
func (c *CatalogEntry30) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry30) String() string { return fmt.Sprintf("CatalogEntry30 %s %s", c.ID, c.Name) }

// CatalogEntry31
type CatalogEntry31 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry31(id,name string) *CatalogEntry31 { return &CatalogEntry31{ID:id,Name:name} }
func (c *CatalogEntry31) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry31) String() string { return fmt.Sprintf("CatalogEntry31 %s %s", c.ID, c.Name) }

// CatalogEntry32
type CatalogEntry32 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry32(id,name string) *CatalogEntry32 { return &CatalogEntry32{ID:id,Name:name} }
func (c *CatalogEntry32) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry32) String() string { return fmt.Sprintf("CatalogEntry32 %s %s", c.ID, c.Name) }

// CatalogEntry33
type CatalogEntry33 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry33(id,name string) *CatalogEntry33 { return &CatalogEntry33{ID:id,Name:name} }
func (c *CatalogEntry33) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry33) String() string { return fmt.Sprintf("CatalogEntry33 %s %s", c.ID, c.Name) }

// CatalogEntry34
type CatalogEntry34 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry34(id,name string) *CatalogEntry34 { return &CatalogEntry34{ID:id,Name:name} }
func (c *CatalogEntry34) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry34) String() string { return fmt.Sprintf("CatalogEntry34 %s %s", c.ID, c.Name) }

// CatalogEntry35
type CatalogEntry35 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry35(id,name string) *CatalogEntry35 { return &CatalogEntry35{ID:id,Name:name} }
func (c *CatalogEntry35) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry35) String() string { return fmt.Sprintf("CatalogEntry35 %s %s", c.ID, c.Name) }

// CatalogEntry36
type CatalogEntry36 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry36(id,name string) *CatalogEntry36 { return &CatalogEntry36{ID:id,Name:name} }
func (c *CatalogEntry36) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry36) String() string { return fmt.Sprintf("CatalogEntry36 %s %s", c.ID, c.Name) }

// CatalogEntry37
type CatalogEntry37 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry37(id,name string) *CatalogEntry37 { return &CatalogEntry37{ID:id,Name:name} }
func (c *CatalogEntry37) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry37) String() string { return fmt.Sprintf("CatalogEntry37 %s %s", c.ID, c.Name) }

// CatalogEntry38
type CatalogEntry38 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry38(id,name string) *CatalogEntry38 { return &CatalogEntry38{ID:id,Name:name} }
func (c *CatalogEntry38) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry38) String() string { return fmt.Sprintf("CatalogEntry38 %s %s", c.ID, c.Name) }

// CatalogEntry39
type CatalogEntry39 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry39(id,name string) *CatalogEntry39 { return &CatalogEntry39{ID:id,Name:name} }
func (c *CatalogEntry39) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry39) String() string { return fmt.Sprintf("CatalogEntry39 %s %s", c.ID, c.Name) }

// CatalogEntry40
type CatalogEntry40 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry40(id,name string) *CatalogEntry40 { return &CatalogEntry40{ID:id,Name:name} }
func (c *CatalogEntry40) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry40) String() string { return fmt.Sprintf("CatalogEntry40 %s %s", c.ID, c.Name) }

// CatalogEntry41
type CatalogEntry41 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry41(id,name string) *CatalogEntry41 { return &CatalogEntry41{ID:id,Name:name} }
func (c *CatalogEntry41) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry41) String() string { return fmt.Sprintf("CatalogEntry41 %s %s", c.ID, c.Name) }

// CatalogEntry42
type CatalogEntry42 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry42(id,name string) *CatalogEntry42 { return &CatalogEntry42{ID:id,Name:name} }
func (c *CatalogEntry42) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry42) String() string { return fmt.Sprintf("CatalogEntry42 %s %s", c.ID, c.Name) }

// CatalogEntry43
type CatalogEntry43 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry43(id,name string) *CatalogEntry43 { return &CatalogEntry43{ID:id,Name:name} }
func (c *CatalogEntry43) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry43) String() string { return fmt.Sprintf("CatalogEntry43 %s %s", c.ID, c.Name) }

// CatalogEntry44
type CatalogEntry44 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry44(id,name string) *CatalogEntry44 { return &CatalogEntry44{ID:id,Name:name} }
func (c *CatalogEntry44) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry44) String() string { return fmt.Sprintf("CatalogEntry44 %s %s", c.ID, c.Name) }

// CatalogEntry45
type CatalogEntry45 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry45(id,name string) *CatalogEntry45 { return &CatalogEntry45{ID:id,Name:name} }
func (c *CatalogEntry45) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry45) String() string { return fmt.Sprintf("CatalogEntry45 %s %s", c.ID, c.Name) }

// CatalogEntry46
type CatalogEntry46 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry46(id,name string) *CatalogEntry46 { return &CatalogEntry46{ID:id,Name:name} }
func (c *CatalogEntry46) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry46) String() string { return fmt.Sprintf("CatalogEntry46 %s %s", c.ID, c.Name) }

// CatalogEntry47
type CatalogEntry47 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry47(id,name string) *CatalogEntry47 { return &CatalogEntry47{ID:id,Name:name} }
func (c *CatalogEntry47) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry47) String() string { return fmt.Sprintf("CatalogEntry47 %s %s", c.ID, c.Name) }

// CatalogEntry48
type CatalogEntry48 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry48(id,name string) *CatalogEntry48 { return &CatalogEntry48{ID:id,Name:name} }
func (c *CatalogEntry48) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry48) String() string { return fmt.Sprintf("CatalogEntry48 %s %s", c.ID, c.Name) }

// CatalogEntry49
type CatalogEntry49 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry49(id,name string) *CatalogEntry49 { return &CatalogEntry49{ID:id,Name:name} }
func (c *CatalogEntry49) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry49) String() string { return fmt.Sprintf("CatalogEntry49 %s %s", c.ID, c.Name) }

// CatalogEntry50
type CatalogEntry50 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry50(id,name string) *CatalogEntry50 { return &CatalogEntry50{ID:id,Name:name} }
func (c *CatalogEntry50) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry50) String() string { return fmt.Sprintf("CatalogEntry50 %s %s", c.ID, c.Name) }

// CatalogEntry51
type CatalogEntry51 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry51(id,name string) *CatalogEntry51 { return &CatalogEntry51{ID:id,Name:name} }
func (c *CatalogEntry51) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry51) String() string { return fmt.Sprintf("CatalogEntry51 %s %s", c.ID, c.Name) }

// CatalogEntry52
type CatalogEntry52 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry52(id,name string) *CatalogEntry52 { return &CatalogEntry52{ID:id,Name:name} }
func (c *CatalogEntry52) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry52) String() string { return fmt.Sprintf("CatalogEntry52 %s %s", c.ID, c.Name) }

// CatalogEntry53
type CatalogEntry53 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry53(id,name string) *CatalogEntry53 { return &CatalogEntry53{ID:id,Name:name} }
func (c *CatalogEntry53) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry53) String() string { return fmt.Sprintf("CatalogEntry53 %s %s", c.ID, c.Name) }

// CatalogEntry54
type CatalogEntry54 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry54(id,name string) *CatalogEntry54 { return &CatalogEntry54{ID:id,Name:name} }
func (c *CatalogEntry54) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry54) String() string { return fmt.Sprintf("CatalogEntry54 %s %s", c.ID, c.Name) }

// CatalogEntry55
type CatalogEntry55 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry55(id,name string) *CatalogEntry55 { return &CatalogEntry55{ID:id,Name:name} }
func (c *CatalogEntry55) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry55) String() string { return fmt.Sprintf("CatalogEntry55 %s %s", c.ID, c.Name) }

// CatalogEntry56
type CatalogEntry56 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry56(id,name string) *CatalogEntry56 { return &CatalogEntry56{ID:id,Name:name} }
func (c *CatalogEntry56) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry56) String() string { return fmt.Sprintf("CatalogEntry56 %s %s", c.ID, c.Name) }

// CatalogEntry57
type CatalogEntry57 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry57(id,name string) *CatalogEntry57 { return &CatalogEntry57{ID:id,Name:name} }
func (c *CatalogEntry57) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry57) String() string { return fmt.Sprintf("CatalogEntry57 %s %s", c.ID, c.Name) }

// CatalogEntry58
type CatalogEntry58 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry58(id,name string) *CatalogEntry58 { return &CatalogEntry58{ID:id,Name:name} }
func (c *CatalogEntry58) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry58) String() string { return fmt.Sprintf("CatalogEntry58 %s %s", c.ID, c.Name) }

// CatalogEntry59
type CatalogEntry59 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry59(id,name string) *CatalogEntry59 { return &CatalogEntry59{ID:id,Name:name} }
func (c *CatalogEntry59) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry59) String() string { return fmt.Sprintf("CatalogEntry59 %s %s", c.ID, c.Name) }

// CatalogEntry60
type CatalogEntry60 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry60(id,name string) *CatalogEntry60 { return &CatalogEntry60{ID:id,Name:name} }
func (c *CatalogEntry60) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry60) String() string { return fmt.Sprintf("CatalogEntry60 %s %s", c.ID, c.Name) }

// CatalogEntry61
type CatalogEntry61 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry61(id,name string) *CatalogEntry61 { return &CatalogEntry61{ID:id,Name:name} }
func (c *CatalogEntry61) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry61) String() string { return fmt.Sprintf("CatalogEntry61 %s %s", c.ID, c.Name) }

// CatalogEntry62
type CatalogEntry62 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry62(id,name string) *CatalogEntry62 { return &CatalogEntry62{ID:id,Name:name} }
func (c *CatalogEntry62) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry62) String() string { return fmt.Sprintf("CatalogEntry62 %s %s", c.ID, c.Name) }

// CatalogEntry63
type CatalogEntry63 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry63(id,name string) *CatalogEntry63 { return &CatalogEntry63{ID:id,Name:name} }
func (c *CatalogEntry63) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry63) String() string { return fmt.Sprintf("CatalogEntry63 %s %s", c.ID, c.Name) }

// CatalogEntry64
type CatalogEntry64 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry64(id,name string) *CatalogEntry64 { return &CatalogEntry64{ID:id,Name:name} }
func (c *CatalogEntry64) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry64) String() string { return fmt.Sprintf("CatalogEntry64 %s %s", c.ID, c.Name) }

// CatalogEntry65
type CatalogEntry65 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry65(id,name string) *CatalogEntry65 { return &CatalogEntry65{ID:id,Name:name} }
func (c *CatalogEntry65) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry65) String() string { return fmt.Sprintf("CatalogEntry65 %s %s", c.ID, c.Name) }

// CatalogEntry66
type CatalogEntry66 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry66(id,name string) *CatalogEntry66 { return &CatalogEntry66{ID:id,Name:name} }
func (c *CatalogEntry66) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry66) String() string { return fmt.Sprintf("CatalogEntry66 %s %s", c.ID, c.Name) }

// CatalogEntry67
type CatalogEntry67 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry67(id,name string) *CatalogEntry67 { return &CatalogEntry67{ID:id,Name:name} }
func (c *CatalogEntry67) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry67) String() string { return fmt.Sprintf("CatalogEntry67 %s %s", c.ID, c.Name) }

// CatalogEntry68
type CatalogEntry68 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry68(id,name string) *CatalogEntry68 { return &CatalogEntry68{ID:id,Name:name} }
func (c *CatalogEntry68) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry68) String() string { return fmt.Sprintf("CatalogEntry68 %s %s", c.ID, c.Name) }

// CatalogEntry69
type CatalogEntry69 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry69(id,name string) *CatalogEntry69 { return &CatalogEntry69{ID:id,Name:name} }
func (c *CatalogEntry69) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry69) String() string { return fmt.Sprintf("CatalogEntry69 %s %s", c.ID, c.Name) }

// CatalogEntry70
type CatalogEntry70 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry70(id,name string) *CatalogEntry70 { return &CatalogEntry70{ID:id,Name:name} }
func (c *CatalogEntry70) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry70) String() string { return fmt.Sprintf("CatalogEntry70 %s %s", c.ID, c.Name) }

// CatalogEntry71
type CatalogEntry71 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry71(id,name string) *CatalogEntry71 { return &CatalogEntry71{ID:id,Name:name} }
func (c *CatalogEntry71) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry71) String() string { return fmt.Sprintf("CatalogEntry71 %s %s", c.ID, c.Name) }

// CatalogEntry72
type CatalogEntry72 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry72(id,name string) *CatalogEntry72 { return &CatalogEntry72{ID:id,Name:name} }
func (c *CatalogEntry72) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry72) String() string { return fmt.Sprintf("CatalogEntry72 %s %s", c.ID, c.Name) }

// CatalogEntry73
type CatalogEntry73 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry73(id,name string) *CatalogEntry73 { return &CatalogEntry73{ID:id,Name:name} }
func (c *CatalogEntry73) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry73) String() string { return fmt.Sprintf("CatalogEntry73 %s %s", c.ID, c.Name) }

// CatalogEntry74
type CatalogEntry74 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry74(id,name string) *CatalogEntry74 { return &CatalogEntry74{ID:id,Name:name} }
func (c *CatalogEntry74) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry74) String() string { return fmt.Sprintf("CatalogEntry74 %s %s", c.ID, c.Name) }

// CatalogEntry75
type CatalogEntry75 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry75(id,name string) *CatalogEntry75 { return &CatalogEntry75{ID:id,Name:name} }
func (c *CatalogEntry75) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry75) String() string { return fmt.Sprintf("CatalogEntry75 %s %s", c.ID, c.Name) }

// CatalogEntry76
type CatalogEntry76 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry76(id,name string) *CatalogEntry76 { return &CatalogEntry76{ID:id,Name:name} }
func (c *CatalogEntry76) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry76) String() string { return fmt.Sprintf("CatalogEntry76 %s %s", c.ID, c.Name) }

// CatalogEntry77
type CatalogEntry77 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry77(id,name string) *CatalogEntry77 { return &CatalogEntry77{ID:id,Name:name} }
func (c *CatalogEntry77) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry77) String() string { return fmt.Sprintf("CatalogEntry77 %s %s", c.ID, c.Name) }

// CatalogEntry78
type CatalogEntry78 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry78(id,name string) *CatalogEntry78 { return &CatalogEntry78{ID:id,Name:name} }
func (c *CatalogEntry78) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry78) String() string { return fmt.Sprintf("CatalogEntry78 %s %s", c.ID, c.Name) }

// CatalogEntry79
type CatalogEntry79 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry79(id,name string) *CatalogEntry79 { return &CatalogEntry79{ID:id,Name:name} }
func (c *CatalogEntry79) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry79) String() string { return fmt.Sprintf("CatalogEntry79 %s %s", c.ID, c.Name) }

// CatalogEntry80
type CatalogEntry80 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry80(id,name string) *CatalogEntry80 { return &CatalogEntry80{ID:id,Name:name} }
func (c *CatalogEntry80) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry80) String() string { return fmt.Sprintf("CatalogEntry80 %s %s", c.ID, c.Name) }

// CatalogEntry81
type CatalogEntry81 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry81(id,name string) *CatalogEntry81 { return &CatalogEntry81{ID:id,Name:name} }
func (c *CatalogEntry81) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry81) String() string { return fmt.Sprintf("CatalogEntry81 %s %s", c.ID, c.Name) }

// CatalogEntry82
type CatalogEntry82 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry82(id,name string) *CatalogEntry82 { return &CatalogEntry82{ID:id,Name:name} }
func (c *CatalogEntry82) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry82) String() string { return fmt.Sprintf("CatalogEntry82 %s %s", c.ID, c.Name) }

// CatalogEntry83
type CatalogEntry83 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry83(id,name string) *CatalogEntry83 { return &CatalogEntry83{ID:id,Name:name} }
func (c *CatalogEntry83) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry83) String() string { return fmt.Sprintf("CatalogEntry83 %s %s", c.ID, c.Name) }

// CatalogEntry84
type CatalogEntry84 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry84(id,name string) *CatalogEntry84 { return &CatalogEntry84{ID:id,Name:name} }
func (c *CatalogEntry84) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry84) String() string { return fmt.Sprintf("CatalogEntry84 %s %s", c.ID, c.Name) }

// CatalogEntry85
type CatalogEntry85 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry85(id,name string) *CatalogEntry85 { return &CatalogEntry85{ID:id,Name:name} }
func (c *CatalogEntry85) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry85) String() string { return fmt.Sprintf("CatalogEntry85 %s %s", c.ID, c.Name) }

// CatalogEntry86
type CatalogEntry86 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry86(id,name string) *CatalogEntry86 { return &CatalogEntry86{ID:id,Name:name} }
func (c *CatalogEntry86) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry86) String() string { return fmt.Sprintf("CatalogEntry86 %s %s", c.ID, c.Name) }

// CatalogEntry87
type CatalogEntry87 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry87(id,name string) *CatalogEntry87 { return &CatalogEntry87{ID:id,Name:name} }
func (c *CatalogEntry87) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry87) String() string { return fmt.Sprintf("CatalogEntry87 %s %s", c.ID, c.Name) }

// CatalogEntry88
type CatalogEntry88 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry88(id,name string) *CatalogEntry88 { return &CatalogEntry88{ID:id,Name:name} }
func (c *CatalogEntry88) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry88) String() string { return fmt.Sprintf("CatalogEntry88 %s %s", c.ID, c.Name) }

// CatalogEntry89
type CatalogEntry89 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry89(id,name string) *CatalogEntry89 { return &CatalogEntry89{ID:id,Name:name} }
func (c *CatalogEntry89) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry89) String() string { return fmt.Sprintf("CatalogEntry89 %s %s", c.ID, c.Name) }

// CatalogEntry90
type CatalogEntry90 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry90(id,name string) *CatalogEntry90 { return &CatalogEntry90{ID:id,Name:name} }
func (c *CatalogEntry90) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry90) String() string { return fmt.Sprintf("CatalogEntry90 %s %s", c.ID, c.Name) }

// CatalogEntry91
type CatalogEntry91 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry91(id,name string) *CatalogEntry91 { return &CatalogEntry91{ID:id,Name:name} }
func (c *CatalogEntry91) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry91) String() string { return fmt.Sprintf("CatalogEntry91 %s %s", c.ID, c.Name) }

// CatalogEntry92
type CatalogEntry92 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry92(id,name string) *CatalogEntry92 { return &CatalogEntry92{ID:id,Name:name} }
func (c *CatalogEntry92) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry92) String() string { return fmt.Sprintf("CatalogEntry92 %s %s", c.ID, c.Name) }

// CatalogEntry93
type CatalogEntry93 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry93(id,name string) *CatalogEntry93 { return &CatalogEntry93{ID:id,Name:name} }
func (c *CatalogEntry93) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry93) String() string { return fmt.Sprintf("CatalogEntry93 %s %s", c.ID, c.Name) }

// CatalogEntry94
type CatalogEntry94 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry94(id,name string) *CatalogEntry94 { return &CatalogEntry94{ID:id,Name:name} }
func (c *CatalogEntry94) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry94) String() string { return fmt.Sprintf("CatalogEntry94 %s %s", c.ID, c.Name) }

// CatalogEntry95
type CatalogEntry95 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry95(id,name string) *CatalogEntry95 { return &CatalogEntry95{ID:id,Name:name} }
func (c *CatalogEntry95) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry95) String() string { return fmt.Sprintf("CatalogEntry95 %s %s", c.ID, c.Name) }

// CatalogEntry96
type CatalogEntry96 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry96(id,name string) *CatalogEntry96 { return &CatalogEntry96{ID:id,Name:name} }
func (c *CatalogEntry96) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry96) String() string { return fmt.Sprintf("CatalogEntry96 %s %s", c.ID, c.Name) }

// CatalogEntry97
type CatalogEntry97 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry97(id,name string) *CatalogEntry97 { return &CatalogEntry97{ID:id,Name:name} }
func (c *CatalogEntry97) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry97) String() string { return fmt.Sprintf("CatalogEntry97 %s %s", c.ID, c.Name) }

// CatalogEntry98
type CatalogEntry98 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry98(id,name string) *CatalogEntry98 { return &CatalogEntry98{ID:id,Name:name} }
func (c *CatalogEntry98) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry98) String() string { return fmt.Sprintf("CatalogEntry98 %s %s", c.ID, c.Name) }

// CatalogEntry99
type CatalogEntry99 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry99(id,name string) *CatalogEntry99 { return &CatalogEntry99{ID:id,Name:name} }
func (c *CatalogEntry99) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry99) String() string { return fmt.Sprintf("CatalogEntry99 %s %s", c.ID, c.Name) }

// CatalogEntry100
type CatalogEntry100 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogEntry100(id,name string) *CatalogEntry100 { return &CatalogEntry100{ID:id,Name:name} }
func (c *CatalogEntry100) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogEntry100) String() string { return fmt.Sprintf("CatalogEntry100 %s %s", c.ID, c.Name) }


