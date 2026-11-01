// Package catalog provides large module 4
package catalog

import (
	"fmt"
	"sync"
)

// CatalogItem1
type CatalogItem1 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem1(id,name string) *CatalogItem1 { return &CatalogItem1{ID:id,Name:name} }
func (c *CatalogItem1) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem1) String() string { return fmt.Sprintf("CatalogItem1 %s %s", c.ID, c.Name) }

// CatalogItem2
type CatalogItem2 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem2(id,name string) *CatalogItem2 { return &CatalogItem2{ID:id,Name:name} }
func (c *CatalogItem2) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem2) String() string { return fmt.Sprintf("CatalogItem2 %s %s", c.ID, c.Name) }

// CatalogItem3
type CatalogItem3 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem3(id,name string) *CatalogItem3 { return &CatalogItem3{ID:id,Name:name} }
func (c *CatalogItem3) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem3) String() string { return fmt.Sprintf("CatalogItem3 %s %s", c.ID, c.Name) }

// CatalogItem4
type CatalogItem4 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem4(id,name string) *CatalogItem4 { return &CatalogItem4{ID:id,Name:name} }
func (c *CatalogItem4) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem4) String() string { return fmt.Sprintf("CatalogItem4 %s %s", c.ID, c.Name) }

// CatalogItem5
type CatalogItem5 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem5(id,name string) *CatalogItem5 { return &CatalogItem5{ID:id,Name:name} }
func (c *CatalogItem5) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem5) String() string { return fmt.Sprintf("CatalogItem5 %s %s", c.ID, c.Name) }

// CatalogItem6
type CatalogItem6 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem6(id,name string) *CatalogItem6 { return &CatalogItem6{ID:id,Name:name} }
func (c *CatalogItem6) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem6) String() string { return fmt.Sprintf("CatalogItem6 %s %s", c.ID, c.Name) }

// CatalogItem7
type CatalogItem7 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem7(id,name string) *CatalogItem7 { return &CatalogItem7{ID:id,Name:name} }
func (c *CatalogItem7) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem7) String() string { return fmt.Sprintf("CatalogItem7 %s %s", c.ID, c.Name) }

// CatalogItem8
type CatalogItem8 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem8(id,name string) *CatalogItem8 { return &CatalogItem8{ID:id,Name:name} }
func (c *CatalogItem8) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem8) String() string { return fmt.Sprintf("CatalogItem8 %s %s", c.ID, c.Name) }

// CatalogItem9
type CatalogItem9 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem9(id,name string) *CatalogItem9 { return &CatalogItem9{ID:id,Name:name} }
func (c *CatalogItem9) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem9) String() string { return fmt.Sprintf("CatalogItem9 %s %s", c.ID, c.Name) }

// CatalogItem10
type CatalogItem10 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem10(id,name string) *CatalogItem10 { return &CatalogItem10{ID:id,Name:name} }
func (c *CatalogItem10) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem10) String() string { return fmt.Sprintf("CatalogItem10 %s %s", c.ID, c.Name) }

// CatalogItem11
type CatalogItem11 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem11(id,name string) *CatalogItem11 { return &CatalogItem11{ID:id,Name:name} }
func (c *CatalogItem11) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem11) String() string { return fmt.Sprintf("CatalogItem11 %s %s", c.ID, c.Name) }

// CatalogItem12
type CatalogItem12 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem12(id,name string) *CatalogItem12 { return &CatalogItem12{ID:id,Name:name} }
func (c *CatalogItem12) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem12) String() string { return fmt.Sprintf("CatalogItem12 %s %s", c.ID, c.Name) }

// CatalogItem13
type CatalogItem13 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem13(id,name string) *CatalogItem13 { return &CatalogItem13{ID:id,Name:name} }
func (c *CatalogItem13) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem13) String() string { return fmt.Sprintf("CatalogItem13 %s %s", c.ID, c.Name) }

// CatalogItem14
type CatalogItem14 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem14(id,name string) *CatalogItem14 { return &CatalogItem14{ID:id,Name:name} }
func (c *CatalogItem14) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem14) String() string { return fmt.Sprintf("CatalogItem14 %s %s", c.ID, c.Name) }

// CatalogItem15
type CatalogItem15 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem15(id,name string) *CatalogItem15 { return &CatalogItem15{ID:id,Name:name} }
func (c *CatalogItem15) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem15) String() string { return fmt.Sprintf("CatalogItem15 %s %s", c.ID, c.Name) }

// CatalogItem16
type CatalogItem16 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem16(id,name string) *CatalogItem16 { return &CatalogItem16{ID:id,Name:name} }
func (c *CatalogItem16) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem16) String() string { return fmt.Sprintf("CatalogItem16 %s %s", c.ID, c.Name) }

// CatalogItem17
type CatalogItem17 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem17(id,name string) *CatalogItem17 { return &CatalogItem17{ID:id,Name:name} }
func (c *CatalogItem17) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem17) String() string { return fmt.Sprintf("CatalogItem17 %s %s", c.ID, c.Name) }

// CatalogItem18
type CatalogItem18 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem18(id,name string) *CatalogItem18 { return &CatalogItem18{ID:id,Name:name} }
func (c *CatalogItem18) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem18) String() string { return fmt.Sprintf("CatalogItem18 %s %s", c.ID, c.Name) }

// CatalogItem19
type CatalogItem19 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem19(id,name string) *CatalogItem19 { return &CatalogItem19{ID:id,Name:name} }
func (c *CatalogItem19) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem19) String() string { return fmt.Sprintf("CatalogItem19 %s %s", c.ID, c.Name) }

// CatalogItem20
type CatalogItem20 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem20(id,name string) *CatalogItem20 { return &CatalogItem20{ID:id,Name:name} }
func (c *CatalogItem20) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem20) String() string { return fmt.Sprintf("CatalogItem20 %s %s", c.ID, c.Name) }

// CatalogItem21
type CatalogItem21 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem21(id,name string) *CatalogItem21 { return &CatalogItem21{ID:id,Name:name} }
func (c *CatalogItem21) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem21) String() string { return fmt.Sprintf("CatalogItem21 %s %s", c.ID, c.Name) }

// CatalogItem22
type CatalogItem22 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem22(id,name string) *CatalogItem22 { return &CatalogItem22{ID:id,Name:name} }
func (c *CatalogItem22) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem22) String() string { return fmt.Sprintf("CatalogItem22 %s %s", c.ID, c.Name) }

// CatalogItem23
type CatalogItem23 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem23(id,name string) *CatalogItem23 { return &CatalogItem23{ID:id,Name:name} }
func (c *CatalogItem23) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem23) String() string { return fmt.Sprintf("CatalogItem23 %s %s", c.ID, c.Name) }

// CatalogItem24
type CatalogItem24 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem24(id,name string) *CatalogItem24 { return &CatalogItem24{ID:id,Name:name} }
func (c *CatalogItem24) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem24) String() string { return fmt.Sprintf("CatalogItem24 %s %s", c.ID, c.Name) }

// CatalogItem25
type CatalogItem25 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem25(id,name string) *CatalogItem25 { return &CatalogItem25{ID:id,Name:name} }
func (c *CatalogItem25) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem25) String() string { return fmt.Sprintf("CatalogItem25 %s %s", c.ID, c.Name) }

// CatalogItem26
type CatalogItem26 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem26(id,name string) *CatalogItem26 { return &CatalogItem26{ID:id,Name:name} }
func (c *CatalogItem26) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem26) String() string { return fmt.Sprintf("CatalogItem26 %s %s", c.ID, c.Name) }

// CatalogItem27
type CatalogItem27 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem27(id,name string) *CatalogItem27 { return &CatalogItem27{ID:id,Name:name} }
func (c *CatalogItem27) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem27) String() string { return fmt.Sprintf("CatalogItem27 %s %s", c.ID, c.Name) }

// CatalogItem28
type CatalogItem28 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem28(id,name string) *CatalogItem28 { return &CatalogItem28{ID:id,Name:name} }
func (c *CatalogItem28) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem28) String() string { return fmt.Sprintf("CatalogItem28 %s %s", c.ID, c.Name) }

// CatalogItem29
type CatalogItem29 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem29(id,name string) *CatalogItem29 { return &CatalogItem29{ID:id,Name:name} }
func (c *CatalogItem29) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem29) String() string { return fmt.Sprintf("CatalogItem29 %s %s", c.ID, c.Name) }

// CatalogItem30
type CatalogItem30 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem30(id,name string) *CatalogItem30 { return &CatalogItem30{ID:id,Name:name} }
func (c *CatalogItem30) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem30) String() string { return fmt.Sprintf("CatalogItem30 %s %s", c.ID, c.Name) }

// CatalogItem31
type CatalogItem31 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem31(id,name string) *CatalogItem31 { return &CatalogItem31{ID:id,Name:name} }
func (c *CatalogItem31) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem31) String() string { return fmt.Sprintf("CatalogItem31 %s %s", c.ID, c.Name) }

// CatalogItem32
type CatalogItem32 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem32(id,name string) *CatalogItem32 { return &CatalogItem32{ID:id,Name:name} }
func (c *CatalogItem32) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem32) String() string { return fmt.Sprintf("CatalogItem32 %s %s", c.ID, c.Name) }

// CatalogItem33
type CatalogItem33 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem33(id,name string) *CatalogItem33 { return &CatalogItem33{ID:id,Name:name} }
func (c *CatalogItem33) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem33) String() string { return fmt.Sprintf("CatalogItem33 %s %s", c.ID, c.Name) }

// CatalogItem34
type CatalogItem34 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem34(id,name string) *CatalogItem34 { return &CatalogItem34{ID:id,Name:name} }
func (c *CatalogItem34) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem34) String() string { return fmt.Sprintf("CatalogItem34 %s %s", c.ID, c.Name) }

// CatalogItem35
type CatalogItem35 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem35(id,name string) *CatalogItem35 { return &CatalogItem35{ID:id,Name:name} }
func (c *CatalogItem35) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem35) String() string { return fmt.Sprintf("CatalogItem35 %s %s", c.ID, c.Name) }

// CatalogItem36
type CatalogItem36 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem36(id,name string) *CatalogItem36 { return &CatalogItem36{ID:id,Name:name} }
func (c *CatalogItem36) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem36) String() string { return fmt.Sprintf("CatalogItem36 %s %s", c.ID, c.Name) }

// CatalogItem37
type CatalogItem37 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem37(id,name string) *CatalogItem37 { return &CatalogItem37{ID:id,Name:name} }
func (c *CatalogItem37) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem37) String() string { return fmt.Sprintf("CatalogItem37 %s %s", c.ID, c.Name) }

// CatalogItem38
type CatalogItem38 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem38(id,name string) *CatalogItem38 { return &CatalogItem38{ID:id,Name:name} }
func (c *CatalogItem38) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem38) String() string { return fmt.Sprintf("CatalogItem38 %s %s", c.ID, c.Name) }

// CatalogItem39
type CatalogItem39 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem39(id,name string) *CatalogItem39 { return &CatalogItem39{ID:id,Name:name} }
func (c *CatalogItem39) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem39) String() string { return fmt.Sprintf("CatalogItem39 %s %s", c.ID, c.Name) }

// CatalogItem40
type CatalogItem40 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem40(id,name string) *CatalogItem40 { return &CatalogItem40{ID:id,Name:name} }
func (c *CatalogItem40) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem40) String() string { return fmt.Sprintf("CatalogItem40 %s %s", c.ID, c.Name) }

// CatalogItem41
type CatalogItem41 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem41(id,name string) *CatalogItem41 { return &CatalogItem41{ID:id,Name:name} }
func (c *CatalogItem41) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem41) String() string { return fmt.Sprintf("CatalogItem41 %s %s", c.ID, c.Name) }

// CatalogItem42
type CatalogItem42 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem42(id,name string) *CatalogItem42 { return &CatalogItem42{ID:id,Name:name} }
func (c *CatalogItem42) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem42) String() string { return fmt.Sprintf("CatalogItem42 %s %s", c.ID, c.Name) }

// CatalogItem43
type CatalogItem43 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem43(id,name string) *CatalogItem43 { return &CatalogItem43{ID:id,Name:name} }
func (c *CatalogItem43) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem43) String() string { return fmt.Sprintf("CatalogItem43 %s %s", c.ID, c.Name) }

// CatalogItem44
type CatalogItem44 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem44(id,name string) *CatalogItem44 { return &CatalogItem44{ID:id,Name:name} }
func (c *CatalogItem44) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem44) String() string { return fmt.Sprintf("CatalogItem44 %s %s", c.ID, c.Name) }

// CatalogItem45
type CatalogItem45 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem45(id,name string) *CatalogItem45 { return &CatalogItem45{ID:id,Name:name} }
func (c *CatalogItem45) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem45) String() string { return fmt.Sprintf("CatalogItem45 %s %s", c.ID, c.Name) }

// CatalogItem46
type CatalogItem46 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem46(id,name string) *CatalogItem46 { return &CatalogItem46{ID:id,Name:name} }
func (c *CatalogItem46) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem46) String() string { return fmt.Sprintf("CatalogItem46 %s %s", c.ID, c.Name) }

// CatalogItem47
type CatalogItem47 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem47(id,name string) *CatalogItem47 { return &CatalogItem47{ID:id,Name:name} }
func (c *CatalogItem47) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem47) String() string { return fmt.Sprintf("CatalogItem47 %s %s", c.ID, c.Name) }

// CatalogItem48
type CatalogItem48 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem48(id,name string) *CatalogItem48 { return &CatalogItem48{ID:id,Name:name} }
func (c *CatalogItem48) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem48) String() string { return fmt.Sprintf("CatalogItem48 %s %s", c.ID, c.Name) }

// CatalogItem49
type CatalogItem49 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem49(id,name string) *CatalogItem49 { return &CatalogItem49{ID:id,Name:name} }
func (c *CatalogItem49) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem49) String() string { return fmt.Sprintf("CatalogItem49 %s %s", c.ID, c.Name) }

// CatalogItem50
type CatalogItem50 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem50(id,name string) *CatalogItem50 { return &CatalogItem50{ID:id,Name:name} }
func (c *CatalogItem50) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem50) String() string { return fmt.Sprintf("CatalogItem50 %s %s", c.ID, c.Name) }

// CatalogItem51
type CatalogItem51 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem51(id,name string) *CatalogItem51 { return &CatalogItem51{ID:id,Name:name} }
func (c *CatalogItem51) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem51) String() string { return fmt.Sprintf("CatalogItem51 %s %s", c.ID, c.Name) }

// CatalogItem52
type CatalogItem52 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem52(id,name string) *CatalogItem52 { return &CatalogItem52{ID:id,Name:name} }
func (c *CatalogItem52) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem52) String() string { return fmt.Sprintf("CatalogItem52 %s %s", c.ID, c.Name) }

// CatalogItem53
type CatalogItem53 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem53(id,name string) *CatalogItem53 { return &CatalogItem53{ID:id,Name:name} }
func (c *CatalogItem53) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem53) String() string { return fmt.Sprintf("CatalogItem53 %s %s", c.ID, c.Name) }

// CatalogItem54
type CatalogItem54 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem54(id,name string) *CatalogItem54 { return &CatalogItem54{ID:id,Name:name} }
func (c *CatalogItem54) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem54) String() string { return fmt.Sprintf("CatalogItem54 %s %s", c.ID, c.Name) }

// CatalogItem55
type CatalogItem55 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem55(id,name string) *CatalogItem55 { return &CatalogItem55{ID:id,Name:name} }
func (c *CatalogItem55) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem55) String() string { return fmt.Sprintf("CatalogItem55 %s %s", c.ID, c.Name) }

// CatalogItem56
type CatalogItem56 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem56(id,name string) *CatalogItem56 { return &CatalogItem56{ID:id,Name:name} }
func (c *CatalogItem56) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem56) String() string { return fmt.Sprintf("CatalogItem56 %s %s", c.ID, c.Name) }

// CatalogItem57
type CatalogItem57 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem57(id,name string) *CatalogItem57 { return &CatalogItem57{ID:id,Name:name} }
func (c *CatalogItem57) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem57) String() string { return fmt.Sprintf("CatalogItem57 %s %s", c.ID, c.Name) }

// CatalogItem58
type CatalogItem58 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem58(id,name string) *CatalogItem58 { return &CatalogItem58{ID:id,Name:name} }
func (c *CatalogItem58) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem58) String() string { return fmt.Sprintf("CatalogItem58 %s %s", c.ID, c.Name) }

// CatalogItem59
type CatalogItem59 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem59(id,name string) *CatalogItem59 { return &CatalogItem59{ID:id,Name:name} }
func (c *CatalogItem59) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem59) String() string { return fmt.Sprintf("CatalogItem59 %s %s", c.ID, c.Name) }

// CatalogItem60
type CatalogItem60 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem60(id,name string) *CatalogItem60 { return &CatalogItem60{ID:id,Name:name} }
func (c *CatalogItem60) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem60) String() string { return fmt.Sprintf("CatalogItem60 %s %s", c.ID, c.Name) }

// CatalogItem61
type CatalogItem61 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem61(id,name string) *CatalogItem61 { return &CatalogItem61{ID:id,Name:name} }
func (c *CatalogItem61) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem61) String() string { return fmt.Sprintf("CatalogItem61 %s %s", c.ID, c.Name) }

// CatalogItem62
type CatalogItem62 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem62(id,name string) *CatalogItem62 { return &CatalogItem62{ID:id,Name:name} }
func (c *CatalogItem62) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem62) String() string { return fmt.Sprintf("CatalogItem62 %s %s", c.ID, c.Name) }

// CatalogItem63
type CatalogItem63 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem63(id,name string) *CatalogItem63 { return &CatalogItem63{ID:id,Name:name} }
func (c *CatalogItem63) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem63) String() string { return fmt.Sprintf("CatalogItem63 %s %s", c.ID, c.Name) }

// CatalogItem64
type CatalogItem64 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem64(id,name string) *CatalogItem64 { return &CatalogItem64{ID:id,Name:name} }
func (c *CatalogItem64) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem64) String() string { return fmt.Sprintf("CatalogItem64 %s %s", c.ID, c.Name) }

// CatalogItem65
type CatalogItem65 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem65(id,name string) *CatalogItem65 { return &CatalogItem65{ID:id,Name:name} }
func (c *CatalogItem65) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem65) String() string { return fmt.Sprintf("CatalogItem65 %s %s", c.ID, c.Name) }

// CatalogItem66
type CatalogItem66 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem66(id,name string) *CatalogItem66 { return &CatalogItem66{ID:id,Name:name} }
func (c *CatalogItem66) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem66) String() string { return fmt.Sprintf("CatalogItem66 %s %s", c.ID, c.Name) }

// CatalogItem67
type CatalogItem67 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem67(id,name string) *CatalogItem67 { return &CatalogItem67{ID:id,Name:name} }
func (c *CatalogItem67) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem67) String() string { return fmt.Sprintf("CatalogItem67 %s %s", c.ID, c.Name) }

// CatalogItem68
type CatalogItem68 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem68(id,name string) *CatalogItem68 { return &CatalogItem68{ID:id,Name:name} }
func (c *CatalogItem68) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem68) String() string { return fmt.Sprintf("CatalogItem68 %s %s", c.ID, c.Name) }

// CatalogItem69
type CatalogItem69 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem69(id,name string) *CatalogItem69 { return &CatalogItem69{ID:id,Name:name} }
func (c *CatalogItem69) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem69) String() string { return fmt.Sprintf("CatalogItem69 %s %s", c.ID, c.Name) }

// CatalogItem70
type CatalogItem70 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem70(id,name string) *CatalogItem70 { return &CatalogItem70{ID:id,Name:name} }
func (c *CatalogItem70) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem70) String() string { return fmt.Sprintf("CatalogItem70 %s %s", c.ID, c.Name) }

// CatalogItem71
type CatalogItem71 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem71(id,name string) *CatalogItem71 { return &CatalogItem71{ID:id,Name:name} }
func (c *CatalogItem71) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem71) String() string { return fmt.Sprintf("CatalogItem71 %s %s", c.ID, c.Name) }

// CatalogItem72
type CatalogItem72 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem72(id,name string) *CatalogItem72 { return &CatalogItem72{ID:id,Name:name} }
func (c *CatalogItem72) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem72) String() string { return fmt.Sprintf("CatalogItem72 %s %s", c.ID, c.Name) }

// CatalogItem73
type CatalogItem73 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem73(id,name string) *CatalogItem73 { return &CatalogItem73{ID:id,Name:name} }
func (c *CatalogItem73) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem73) String() string { return fmt.Sprintf("CatalogItem73 %s %s", c.ID, c.Name) }

// CatalogItem74
type CatalogItem74 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem74(id,name string) *CatalogItem74 { return &CatalogItem74{ID:id,Name:name} }
func (c *CatalogItem74) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem74) String() string { return fmt.Sprintf("CatalogItem74 %s %s", c.ID, c.Name) }

// CatalogItem75
type CatalogItem75 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem75(id,name string) *CatalogItem75 { return &CatalogItem75{ID:id,Name:name} }
func (c *CatalogItem75) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem75) String() string { return fmt.Sprintf("CatalogItem75 %s %s", c.ID, c.Name) }

// CatalogItem76
type CatalogItem76 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem76(id,name string) *CatalogItem76 { return &CatalogItem76{ID:id,Name:name} }
func (c *CatalogItem76) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem76) String() string { return fmt.Sprintf("CatalogItem76 %s %s", c.ID, c.Name) }

// CatalogItem77
type CatalogItem77 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem77(id,name string) *CatalogItem77 { return &CatalogItem77{ID:id,Name:name} }
func (c *CatalogItem77) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem77) String() string { return fmt.Sprintf("CatalogItem77 %s %s", c.ID, c.Name) }

// CatalogItem78
type CatalogItem78 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem78(id,name string) *CatalogItem78 { return &CatalogItem78{ID:id,Name:name} }
func (c *CatalogItem78) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem78) String() string { return fmt.Sprintf("CatalogItem78 %s %s", c.ID, c.Name) }

// CatalogItem79
type CatalogItem79 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem79(id,name string) *CatalogItem79 { return &CatalogItem79{ID:id,Name:name} }
func (c *CatalogItem79) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem79) String() string { return fmt.Sprintf("CatalogItem79 %s %s", c.ID, c.Name) }

// CatalogItem80
type CatalogItem80 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogItem80(id,name string) *CatalogItem80 { return &CatalogItem80{ID:id,Name:name} }
func (c *CatalogItem80) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogItem80) String() string { return fmt.Sprintf("CatalogItem80 %s %s", c.ID, c.Name) }


