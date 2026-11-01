// Package catalog provides large module 8
package catalog

import (
	"fmt"
	"sync"
)

// CatalogNode8_1
type CatalogNode8_1 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_1(id,name string) *CatalogNode8_1 { return &CatalogNode8_1{ID:id,Name:name} }
func (c *CatalogNode8_1) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_1) String() string { return fmt.Sprintf("CatalogNode8_1 %s %s", c.ID, c.Name) }

// CatalogNode8_2
type CatalogNode8_2 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_2(id,name string) *CatalogNode8_2 { return &CatalogNode8_2{ID:id,Name:name} }
func (c *CatalogNode8_2) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_2) String() string { return fmt.Sprintf("CatalogNode8_2 %s %s", c.ID, c.Name) }

// CatalogNode8_3
type CatalogNode8_3 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_3(id,name string) *CatalogNode8_3 { return &CatalogNode8_3{ID:id,Name:name} }
func (c *CatalogNode8_3) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_3) String() string { return fmt.Sprintf("CatalogNode8_3 %s %s", c.ID, c.Name) }

// CatalogNode8_4
type CatalogNode8_4 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_4(id,name string) *CatalogNode8_4 { return &CatalogNode8_4{ID:id,Name:name} }
func (c *CatalogNode8_4) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_4) String() string { return fmt.Sprintf("CatalogNode8_4 %s %s", c.ID, c.Name) }

// CatalogNode8_5
type CatalogNode8_5 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_5(id,name string) *CatalogNode8_5 { return &CatalogNode8_5{ID:id,Name:name} }
func (c *CatalogNode8_5) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_5) String() string { return fmt.Sprintf("CatalogNode8_5 %s %s", c.ID, c.Name) }

// CatalogNode8_6
type CatalogNode8_6 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_6(id,name string) *CatalogNode8_6 { return &CatalogNode8_6{ID:id,Name:name} }
func (c *CatalogNode8_6) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_6) String() string { return fmt.Sprintf("CatalogNode8_6 %s %s", c.ID, c.Name) }

// CatalogNode8_7
type CatalogNode8_7 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_7(id,name string) *CatalogNode8_7 { return &CatalogNode8_7{ID:id,Name:name} }
func (c *CatalogNode8_7) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_7) String() string { return fmt.Sprintf("CatalogNode8_7 %s %s", c.ID, c.Name) }

// CatalogNode8_8
type CatalogNode8_8 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_8(id,name string) *CatalogNode8_8 { return &CatalogNode8_8{ID:id,Name:name} }
func (c *CatalogNode8_8) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_8) String() string { return fmt.Sprintf("CatalogNode8_8 %s %s", c.ID, c.Name) }

// CatalogNode8_9
type CatalogNode8_9 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_9(id,name string) *CatalogNode8_9 { return &CatalogNode8_9{ID:id,Name:name} }
func (c *CatalogNode8_9) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_9) String() string { return fmt.Sprintf("CatalogNode8_9 %s %s", c.ID, c.Name) }

// CatalogNode8_10
type CatalogNode8_10 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_10(id,name string) *CatalogNode8_10 { return &CatalogNode8_10{ID:id,Name:name} }
func (c *CatalogNode8_10) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_10) String() string { return fmt.Sprintf("CatalogNode8_10 %s %s", c.ID, c.Name) }

// CatalogNode8_11
type CatalogNode8_11 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_11(id,name string) *CatalogNode8_11 { return &CatalogNode8_11{ID:id,Name:name} }
func (c *CatalogNode8_11) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_11) String() string { return fmt.Sprintf("CatalogNode8_11 %s %s", c.ID, c.Name) }

// CatalogNode8_12
type CatalogNode8_12 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_12(id,name string) *CatalogNode8_12 { return &CatalogNode8_12{ID:id,Name:name} }
func (c *CatalogNode8_12) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_12) String() string { return fmt.Sprintf("CatalogNode8_12 %s %s", c.ID, c.Name) }

// CatalogNode8_13
type CatalogNode8_13 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_13(id,name string) *CatalogNode8_13 { return &CatalogNode8_13{ID:id,Name:name} }
func (c *CatalogNode8_13) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_13) String() string { return fmt.Sprintf("CatalogNode8_13 %s %s", c.ID, c.Name) }

// CatalogNode8_14
type CatalogNode8_14 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_14(id,name string) *CatalogNode8_14 { return &CatalogNode8_14{ID:id,Name:name} }
func (c *CatalogNode8_14) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_14) String() string { return fmt.Sprintf("CatalogNode8_14 %s %s", c.ID, c.Name) }

// CatalogNode8_15
type CatalogNode8_15 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_15(id,name string) *CatalogNode8_15 { return &CatalogNode8_15{ID:id,Name:name} }
func (c *CatalogNode8_15) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_15) String() string { return fmt.Sprintf("CatalogNode8_15 %s %s", c.ID, c.Name) }

// CatalogNode8_16
type CatalogNode8_16 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_16(id,name string) *CatalogNode8_16 { return &CatalogNode8_16{ID:id,Name:name} }
func (c *CatalogNode8_16) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_16) String() string { return fmt.Sprintf("CatalogNode8_16 %s %s", c.ID, c.Name) }

// CatalogNode8_17
type CatalogNode8_17 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_17(id,name string) *CatalogNode8_17 { return &CatalogNode8_17{ID:id,Name:name} }
func (c *CatalogNode8_17) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_17) String() string { return fmt.Sprintf("CatalogNode8_17 %s %s", c.ID, c.Name) }

// CatalogNode8_18
type CatalogNode8_18 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_18(id,name string) *CatalogNode8_18 { return &CatalogNode8_18{ID:id,Name:name} }
func (c *CatalogNode8_18) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_18) String() string { return fmt.Sprintf("CatalogNode8_18 %s %s", c.ID, c.Name) }

// CatalogNode8_19
type CatalogNode8_19 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_19(id,name string) *CatalogNode8_19 { return &CatalogNode8_19{ID:id,Name:name} }
func (c *CatalogNode8_19) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_19) String() string { return fmt.Sprintf("CatalogNode8_19 %s %s", c.ID, c.Name) }

// CatalogNode8_20
type CatalogNode8_20 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_20(id,name string) *CatalogNode8_20 { return &CatalogNode8_20{ID:id,Name:name} }
func (c *CatalogNode8_20) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_20) String() string { return fmt.Sprintf("CatalogNode8_20 %s %s", c.ID, c.Name) }

// CatalogNode8_21
type CatalogNode8_21 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_21(id,name string) *CatalogNode8_21 { return &CatalogNode8_21{ID:id,Name:name} }
func (c *CatalogNode8_21) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_21) String() string { return fmt.Sprintf("CatalogNode8_21 %s %s", c.ID, c.Name) }

// CatalogNode8_22
type CatalogNode8_22 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_22(id,name string) *CatalogNode8_22 { return &CatalogNode8_22{ID:id,Name:name} }
func (c *CatalogNode8_22) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_22) String() string { return fmt.Sprintf("CatalogNode8_22 %s %s", c.ID, c.Name) }

// CatalogNode8_23
type CatalogNode8_23 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_23(id,name string) *CatalogNode8_23 { return &CatalogNode8_23{ID:id,Name:name} }
func (c *CatalogNode8_23) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_23) String() string { return fmt.Sprintf("CatalogNode8_23 %s %s", c.ID, c.Name) }

// CatalogNode8_24
type CatalogNode8_24 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_24(id,name string) *CatalogNode8_24 { return &CatalogNode8_24{ID:id,Name:name} }
func (c *CatalogNode8_24) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_24) String() string { return fmt.Sprintf("CatalogNode8_24 %s %s", c.ID, c.Name) }

// CatalogNode8_25
type CatalogNode8_25 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_25(id,name string) *CatalogNode8_25 { return &CatalogNode8_25{ID:id,Name:name} }
func (c *CatalogNode8_25) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_25) String() string { return fmt.Sprintf("CatalogNode8_25 %s %s", c.ID, c.Name) }

// CatalogNode8_26
type CatalogNode8_26 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_26(id,name string) *CatalogNode8_26 { return &CatalogNode8_26{ID:id,Name:name} }
func (c *CatalogNode8_26) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_26) String() string { return fmt.Sprintf("CatalogNode8_26 %s %s", c.ID, c.Name) }

// CatalogNode8_27
type CatalogNode8_27 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_27(id,name string) *CatalogNode8_27 { return &CatalogNode8_27{ID:id,Name:name} }
func (c *CatalogNode8_27) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_27) String() string { return fmt.Sprintf("CatalogNode8_27 %s %s", c.ID, c.Name) }

// CatalogNode8_28
type CatalogNode8_28 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_28(id,name string) *CatalogNode8_28 { return &CatalogNode8_28{ID:id,Name:name} }
func (c *CatalogNode8_28) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_28) String() string { return fmt.Sprintf("CatalogNode8_28 %s %s", c.ID, c.Name) }

// CatalogNode8_29
type CatalogNode8_29 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_29(id,name string) *CatalogNode8_29 { return &CatalogNode8_29{ID:id,Name:name} }
func (c *CatalogNode8_29) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_29) String() string { return fmt.Sprintf("CatalogNode8_29 %s %s", c.ID, c.Name) }

// CatalogNode8_30
type CatalogNode8_30 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_30(id,name string) *CatalogNode8_30 { return &CatalogNode8_30{ID:id,Name:name} }
func (c *CatalogNode8_30) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_30) String() string { return fmt.Sprintf("CatalogNode8_30 %s %s", c.ID, c.Name) }

// CatalogNode8_31
type CatalogNode8_31 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_31(id,name string) *CatalogNode8_31 { return &CatalogNode8_31{ID:id,Name:name} }
func (c *CatalogNode8_31) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_31) String() string { return fmt.Sprintf("CatalogNode8_31 %s %s", c.ID, c.Name) }

// CatalogNode8_32
type CatalogNode8_32 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_32(id,name string) *CatalogNode8_32 { return &CatalogNode8_32{ID:id,Name:name} }
func (c *CatalogNode8_32) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_32) String() string { return fmt.Sprintf("CatalogNode8_32 %s %s", c.ID, c.Name) }

// CatalogNode8_33
type CatalogNode8_33 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_33(id,name string) *CatalogNode8_33 { return &CatalogNode8_33{ID:id,Name:name} }
func (c *CatalogNode8_33) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_33) String() string { return fmt.Sprintf("CatalogNode8_33 %s %s", c.ID, c.Name) }

// CatalogNode8_34
type CatalogNode8_34 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_34(id,name string) *CatalogNode8_34 { return &CatalogNode8_34{ID:id,Name:name} }
func (c *CatalogNode8_34) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_34) String() string { return fmt.Sprintf("CatalogNode8_34 %s %s", c.ID, c.Name) }

// CatalogNode8_35
type CatalogNode8_35 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_35(id,name string) *CatalogNode8_35 { return &CatalogNode8_35{ID:id,Name:name} }
func (c *CatalogNode8_35) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_35) String() string { return fmt.Sprintf("CatalogNode8_35 %s %s", c.ID, c.Name) }

// CatalogNode8_36
type CatalogNode8_36 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_36(id,name string) *CatalogNode8_36 { return &CatalogNode8_36{ID:id,Name:name} }
func (c *CatalogNode8_36) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_36) String() string { return fmt.Sprintf("CatalogNode8_36 %s %s", c.ID, c.Name) }

// CatalogNode8_37
type CatalogNode8_37 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_37(id,name string) *CatalogNode8_37 { return &CatalogNode8_37{ID:id,Name:name} }
func (c *CatalogNode8_37) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_37) String() string { return fmt.Sprintf("CatalogNode8_37 %s %s", c.ID, c.Name) }

// CatalogNode8_38
type CatalogNode8_38 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_38(id,name string) *CatalogNode8_38 { return &CatalogNode8_38{ID:id,Name:name} }
func (c *CatalogNode8_38) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_38) String() string { return fmt.Sprintf("CatalogNode8_38 %s %s", c.ID, c.Name) }

// CatalogNode8_39
type CatalogNode8_39 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_39(id,name string) *CatalogNode8_39 { return &CatalogNode8_39{ID:id,Name:name} }
func (c *CatalogNode8_39) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_39) String() string { return fmt.Sprintf("CatalogNode8_39 %s %s", c.ID, c.Name) }

// CatalogNode8_40
type CatalogNode8_40 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_40(id,name string) *CatalogNode8_40 { return &CatalogNode8_40{ID:id,Name:name} }
func (c *CatalogNode8_40) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_40) String() string { return fmt.Sprintf("CatalogNode8_40 %s %s", c.ID, c.Name) }

// CatalogNode8_41
type CatalogNode8_41 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_41(id,name string) *CatalogNode8_41 { return &CatalogNode8_41{ID:id,Name:name} }
func (c *CatalogNode8_41) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_41) String() string { return fmt.Sprintf("CatalogNode8_41 %s %s", c.ID, c.Name) }

// CatalogNode8_42
type CatalogNode8_42 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_42(id,name string) *CatalogNode8_42 { return &CatalogNode8_42{ID:id,Name:name} }
func (c *CatalogNode8_42) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_42) String() string { return fmt.Sprintf("CatalogNode8_42 %s %s", c.ID, c.Name) }

// CatalogNode8_43
type CatalogNode8_43 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_43(id,name string) *CatalogNode8_43 { return &CatalogNode8_43{ID:id,Name:name} }
func (c *CatalogNode8_43) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_43) String() string { return fmt.Sprintf("CatalogNode8_43 %s %s", c.ID, c.Name) }

// CatalogNode8_44
type CatalogNode8_44 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_44(id,name string) *CatalogNode8_44 { return &CatalogNode8_44{ID:id,Name:name} }
func (c *CatalogNode8_44) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_44) String() string { return fmt.Sprintf("CatalogNode8_44 %s %s", c.ID, c.Name) }

// CatalogNode8_45
type CatalogNode8_45 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_45(id,name string) *CatalogNode8_45 { return &CatalogNode8_45{ID:id,Name:name} }
func (c *CatalogNode8_45) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_45) String() string { return fmt.Sprintf("CatalogNode8_45 %s %s", c.ID, c.Name) }

// CatalogNode8_46
type CatalogNode8_46 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_46(id,name string) *CatalogNode8_46 { return &CatalogNode8_46{ID:id,Name:name} }
func (c *CatalogNode8_46) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_46) String() string { return fmt.Sprintf("CatalogNode8_46 %s %s", c.ID, c.Name) }

// CatalogNode8_47
type CatalogNode8_47 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_47(id,name string) *CatalogNode8_47 { return &CatalogNode8_47{ID:id,Name:name} }
func (c *CatalogNode8_47) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_47) String() string { return fmt.Sprintf("CatalogNode8_47 %s %s", c.ID, c.Name) }

// CatalogNode8_48
type CatalogNode8_48 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_48(id,name string) *CatalogNode8_48 { return &CatalogNode8_48{ID:id,Name:name} }
func (c *CatalogNode8_48) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_48) String() string { return fmt.Sprintf("CatalogNode8_48 %s %s", c.ID, c.Name) }

// CatalogNode8_49
type CatalogNode8_49 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_49(id,name string) *CatalogNode8_49 { return &CatalogNode8_49{ID:id,Name:name} }
func (c *CatalogNode8_49) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_49) String() string { return fmt.Sprintf("CatalogNode8_49 %s %s", c.ID, c.Name) }

// CatalogNode8_50
type CatalogNode8_50 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_50(id,name string) *CatalogNode8_50 { return &CatalogNode8_50{ID:id,Name:name} }
func (c *CatalogNode8_50) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_50) String() string { return fmt.Sprintf("CatalogNode8_50 %s %s", c.ID, c.Name) }

// CatalogNode8_51
type CatalogNode8_51 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_51(id,name string) *CatalogNode8_51 { return &CatalogNode8_51{ID:id,Name:name} }
func (c *CatalogNode8_51) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_51) String() string { return fmt.Sprintf("CatalogNode8_51 %s %s", c.ID, c.Name) }

// CatalogNode8_52
type CatalogNode8_52 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_52(id,name string) *CatalogNode8_52 { return &CatalogNode8_52{ID:id,Name:name} }
func (c *CatalogNode8_52) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_52) String() string { return fmt.Sprintf("CatalogNode8_52 %s %s", c.ID, c.Name) }

// CatalogNode8_53
type CatalogNode8_53 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_53(id,name string) *CatalogNode8_53 { return &CatalogNode8_53{ID:id,Name:name} }
func (c *CatalogNode8_53) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_53) String() string { return fmt.Sprintf("CatalogNode8_53 %s %s", c.ID, c.Name) }

// CatalogNode8_54
type CatalogNode8_54 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_54(id,name string) *CatalogNode8_54 { return &CatalogNode8_54{ID:id,Name:name} }
func (c *CatalogNode8_54) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_54) String() string { return fmt.Sprintf("CatalogNode8_54 %s %s", c.ID, c.Name) }

// CatalogNode8_55
type CatalogNode8_55 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_55(id,name string) *CatalogNode8_55 { return &CatalogNode8_55{ID:id,Name:name} }
func (c *CatalogNode8_55) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_55) String() string { return fmt.Sprintf("CatalogNode8_55 %s %s", c.ID, c.Name) }

// CatalogNode8_56
type CatalogNode8_56 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_56(id,name string) *CatalogNode8_56 { return &CatalogNode8_56{ID:id,Name:name} }
func (c *CatalogNode8_56) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_56) String() string { return fmt.Sprintf("CatalogNode8_56 %s %s", c.ID, c.Name) }

// CatalogNode8_57
type CatalogNode8_57 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_57(id,name string) *CatalogNode8_57 { return &CatalogNode8_57{ID:id,Name:name} }
func (c *CatalogNode8_57) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_57) String() string { return fmt.Sprintf("CatalogNode8_57 %s %s", c.ID, c.Name) }

// CatalogNode8_58
type CatalogNode8_58 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_58(id,name string) *CatalogNode8_58 { return &CatalogNode8_58{ID:id,Name:name} }
func (c *CatalogNode8_58) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_58) String() string { return fmt.Sprintf("CatalogNode8_58 %s %s", c.ID, c.Name) }

// CatalogNode8_59
type CatalogNode8_59 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_59(id,name string) *CatalogNode8_59 { return &CatalogNode8_59{ID:id,Name:name} }
func (c *CatalogNode8_59) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_59) String() string { return fmt.Sprintf("CatalogNode8_59 %s %s", c.ID, c.Name) }

// CatalogNode8_60
type CatalogNode8_60 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_60(id,name string) *CatalogNode8_60 { return &CatalogNode8_60{ID:id,Name:name} }
func (c *CatalogNode8_60) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_60) String() string { return fmt.Sprintf("CatalogNode8_60 %s %s", c.ID, c.Name) }

// CatalogNode8_61
type CatalogNode8_61 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_61(id,name string) *CatalogNode8_61 { return &CatalogNode8_61{ID:id,Name:name} }
func (c *CatalogNode8_61) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_61) String() string { return fmt.Sprintf("CatalogNode8_61 %s %s", c.ID, c.Name) }

// CatalogNode8_62
type CatalogNode8_62 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_62(id,name string) *CatalogNode8_62 { return &CatalogNode8_62{ID:id,Name:name} }
func (c *CatalogNode8_62) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_62) String() string { return fmt.Sprintf("CatalogNode8_62 %s %s", c.ID, c.Name) }

// CatalogNode8_63
type CatalogNode8_63 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_63(id,name string) *CatalogNode8_63 { return &CatalogNode8_63{ID:id,Name:name} }
func (c *CatalogNode8_63) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_63) String() string { return fmt.Sprintf("CatalogNode8_63 %s %s", c.ID, c.Name) }

// CatalogNode8_64
type CatalogNode8_64 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_64(id,name string) *CatalogNode8_64 { return &CatalogNode8_64{ID:id,Name:name} }
func (c *CatalogNode8_64) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_64) String() string { return fmt.Sprintf("CatalogNode8_64 %s %s", c.ID, c.Name) }

// CatalogNode8_65
type CatalogNode8_65 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_65(id,name string) *CatalogNode8_65 { return &CatalogNode8_65{ID:id,Name:name} }
func (c *CatalogNode8_65) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_65) String() string { return fmt.Sprintf("CatalogNode8_65 %s %s", c.ID, c.Name) }

// CatalogNode8_66
type CatalogNode8_66 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_66(id,name string) *CatalogNode8_66 { return &CatalogNode8_66{ID:id,Name:name} }
func (c *CatalogNode8_66) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_66) String() string { return fmt.Sprintf("CatalogNode8_66 %s %s", c.ID, c.Name) }

// CatalogNode8_67
type CatalogNode8_67 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_67(id,name string) *CatalogNode8_67 { return &CatalogNode8_67{ID:id,Name:name} }
func (c *CatalogNode8_67) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_67) String() string { return fmt.Sprintf("CatalogNode8_67 %s %s", c.ID, c.Name) }

// CatalogNode8_68
type CatalogNode8_68 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_68(id,name string) *CatalogNode8_68 { return &CatalogNode8_68{ID:id,Name:name} }
func (c *CatalogNode8_68) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_68) String() string { return fmt.Sprintf("CatalogNode8_68 %s %s", c.ID, c.Name) }

// CatalogNode8_69
type CatalogNode8_69 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_69(id,name string) *CatalogNode8_69 { return &CatalogNode8_69{ID:id,Name:name} }
func (c *CatalogNode8_69) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_69) String() string { return fmt.Sprintf("CatalogNode8_69 %s %s", c.ID, c.Name) }

// CatalogNode8_70
type CatalogNode8_70 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_70(id,name string) *CatalogNode8_70 { return &CatalogNode8_70{ID:id,Name:name} }
func (c *CatalogNode8_70) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_70) String() string { return fmt.Sprintf("CatalogNode8_70 %s %s", c.ID, c.Name) }

// CatalogNode8_71
type CatalogNode8_71 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_71(id,name string) *CatalogNode8_71 { return &CatalogNode8_71{ID:id,Name:name} }
func (c *CatalogNode8_71) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_71) String() string { return fmt.Sprintf("CatalogNode8_71 %s %s", c.ID, c.Name) }

// CatalogNode8_72
type CatalogNode8_72 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_72(id,name string) *CatalogNode8_72 { return &CatalogNode8_72{ID:id,Name:name} }
func (c *CatalogNode8_72) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_72) String() string { return fmt.Sprintf("CatalogNode8_72 %s %s", c.ID, c.Name) }

// CatalogNode8_73
type CatalogNode8_73 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_73(id,name string) *CatalogNode8_73 { return &CatalogNode8_73{ID:id,Name:name} }
func (c *CatalogNode8_73) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_73) String() string { return fmt.Sprintf("CatalogNode8_73 %s %s", c.ID, c.Name) }

// CatalogNode8_74
type CatalogNode8_74 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_74(id,name string) *CatalogNode8_74 { return &CatalogNode8_74{ID:id,Name:name} }
func (c *CatalogNode8_74) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_74) String() string { return fmt.Sprintf("CatalogNode8_74 %s %s", c.ID, c.Name) }

// CatalogNode8_75
type CatalogNode8_75 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_75(id,name string) *CatalogNode8_75 { return &CatalogNode8_75{ID:id,Name:name} }
func (c *CatalogNode8_75) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_75) String() string { return fmt.Sprintf("CatalogNode8_75 %s %s", c.ID, c.Name) }

// CatalogNode8_76
type CatalogNode8_76 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_76(id,name string) *CatalogNode8_76 { return &CatalogNode8_76{ID:id,Name:name} }
func (c *CatalogNode8_76) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_76) String() string { return fmt.Sprintf("CatalogNode8_76 %s %s", c.ID, c.Name) }

// CatalogNode8_77
type CatalogNode8_77 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_77(id,name string) *CatalogNode8_77 { return &CatalogNode8_77{ID:id,Name:name} }
func (c *CatalogNode8_77) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_77) String() string { return fmt.Sprintf("CatalogNode8_77 %s %s", c.ID, c.Name) }

// CatalogNode8_78
type CatalogNode8_78 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_78(id,name string) *CatalogNode8_78 { return &CatalogNode8_78{ID:id,Name:name} }
func (c *CatalogNode8_78) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_78) String() string { return fmt.Sprintf("CatalogNode8_78 %s %s", c.ID, c.Name) }

// CatalogNode8_79
type CatalogNode8_79 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_79(id,name string) *CatalogNode8_79 { return &CatalogNode8_79{ID:id,Name:name} }
func (c *CatalogNode8_79) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_79) String() string { return fmt.Sprintf("CatalogNode8_79 %s %s", c.ID, c.Name) }

// CatalogNode8_80
type CatalogNode8_80 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_80(id,name string) *CatalogNode8_80 { return &CatalogNode8_80{ID:id,Name:name} }
func (c *CatalogNode8_80) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_80) String() string { return fmt.Sprintf("CatalogNode8_80 %s %s", c.ID, c.Name) }

// CatalogNode8_81
type CatalogNode8_81 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_81(id,name string) *CatalogNode8_81 { return &CatalogNode8_81{ID:id,Name:name} }
func (c *CatalogNode8_81) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_81) String() string { return fmt.Sprintf("CatalogNode8_81 %s %s", c.ID, c.Name) }

// CatalogNode8_82
type CatalogNode8_82 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_82(id,name string) *CatalogNode8_82 { return &CatalogNode8_82{ID:id,Name:name} }
func (c *CatalogNode8_82) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_82) String() string { return fmt.Sprintf("CatalogNode8_82 %s %s", c.ID, c.Name) }

// CatalogNode8_83
type CatalogNode8_83 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_83(id,name string) *CatalogNode8_83 { return &CatalogNode8_83{ID:id,Name:name} }
func (c *CatalogNode8_83) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_83) String() string { return fmt.Sprintf("CatalogNode8_83 %s %s", c.ID, c.Name) }

// CatalogNode8_84
type CatalogNode8_84 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_84(id,name string) *CatalogNode8_84 { return &CatalogNode8_84{ID:id,Name:name} }
func (c *CatalogNode8_84) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_84) String() string { return fmt.Sprintf("CatalogNode8_84 %s %s", c.ID, c.Name) }

// CatalogNode8_85
type CatalogNode8_85 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_85(id,name string) *CatalogNode8_85 { return &CatalogNode8_85{ID:id,Name:name} }
func (c *CatalogNode8_85) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_85) String() string { return fmt.Sprintf("CatalogNode8_85 %s %s", c.ID, c.Name) }

// CatalogNode8_86
type CatalogNode8_86 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_86(id,name string) *CatalogNode8_86 { return &CatalogNode8_86{ID:id,Name:name} }
func (c *CatalogNode8_86) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_86) String() string { return fmt.Sprintf("CatalogNode8_86 %s %s", c.ID, c.Name) }

// CatalogNode8_87
type CatalogNode8_87 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_87(id,name string) *CatalogNode8_87 { return &CatalogNode8_87{ID:id,Name:name} }
func (c *CatalogNode8_87) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_87) String() string { return fmt.Sprintf("CatalogNode8_87 %s %s", c.ID, c.Name) }

// CatalogNode8_88
type CatalogNode8_88 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_88(id,name string) *CatalogNode8_88 { return &CatalogNode8_88{ID:id,Name:name} }
func (c *CatalogNode8_88) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_88) String() string { return fmt.Sprintf("CatalogNode8_88 %s %s", c.ID, c.Name) }

// CatalogNode8_89
type CatalogNode8_89 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_89(id,name string) *CatalogNode8_89 { return &CatalogNode8_89{ID:id,Name:name} }
func (c *CatalogNode8_89) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_89) String() string { return fmt.Sprintf("CatalogNode8_89 %s %s", c.ID, c.Name) }

// CatalogNode8_90
type CatalogNode8_90 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_90(id,name string) *CatalogNode8_90 { return &CatalogNode8_90{ID:id,Name:name} }
func (c *CatalogNode8_90) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_90) String() string { return fmt.Sprintf("CatalogNode8_90 %s %s", c.ID, c.Name) }

// CatalogNode8_91
type CatalogNode8_91 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_91(id,name string) *CatalogNode8_91 { return &CatalogNode8_91{ID:id,Name:name} }
func (c *CatalogNode8_91) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_91) String() string { return fmt.Sprintf("CatalogNode8_91 %s %s", c.ID, c.Name) }

// CatalogNode8_92
type CatalogNode8_92 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_92(id,name string) *CatalogNode8_92 { return &CatalogNode8_92{ID:id,Name:name} }
func (c *CatalogNode8_92) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_92) String() string { return fmt.Sprintf("CatalogNode8_92 %s %s", c.ID, c.Name) }

// CatalogNode8_93
type CatalogNode8_93 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_93(id,name string) *CatalogNode8_93 { return &CatalogNode8_93{ID:id,Name:name} }
func (c *CatalogNode8_93) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_93) String() string { return fmt.Sprintf("CatalogNode8_93 %s %s", c.ID, c.Name) }

// CatalogNode8_94
type CatalogNode8_94 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_94(id,name string) *CatalogNode8_94 { return &CatalogNode8_94{ID:id,Name:name} }
func (c *CatalogNode8_94) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_94) String() string { return fmt.Sprintf("CatalogNode8_94 %s %s", c.ID, c.Name) }

// CatalogNode8_95
type CatalogNode8_95 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_95(id,name string) *CatalogNode8_95 { return &CatalogNode8_95{ID:id,Name:name} }
func (c *CatalogNode8_95) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_95) String() string { return fmt.Sprintf("CatalogNode8_95 %s %s", c.ID, c.Name) }

// CatalogNode8_96
type CatalogNode8_96 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_96(id,name string) *CatalogNode8_96 { return &CatalogNode8_96{ID:id,Name:name} }
func (c *CatalogNode8_96) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_96) String() string { return fmt.Sprintf("CatalogNode8_96 %s %s", c.ID, c.Name) }

// CatalogNode8_97
type CatalogNode8_97 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_97(id,name string) *CatalogNode8_97 { return &CatalogNode8_97{ID:id,Name:name} }
func (c *CatalogNode8_97) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_97) String() string { return fmt.Sprintf("CatalogNode8_97 %s %s", c.ID, c.Name) }

// CatalogNode8_98
type CatalogNode8_98 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_98(id,name string) *CatalogNode8_98 { return &CatalogNode8_98{ID:id,Name:name} }
func (c *CatalogNode8_98) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_98) String() string { return fmt.Sprintf("CatalogNode8_98 %s %s", c.ID, c.Name) }

// CatalogNode8_99
type CatalogNode8_99 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_99(id,name string) *CatalogNode8_99 { return &CatalogNode8_99{ID:id,Name:name} }
func (c *CatalogNode8_99) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_99) String() string { return fmt.Sprintf("CatalogNode8_99 %s %s", c.ID, c.Name) }

// CatalogNode8_100
type CatalogNode8_100 struct {
	ID string
	Name string
	mu sync.Mutex
}

func NewCatalogNode8_100(id,name string) *CatalogNode8_100 { return &CatalogNode8_100{ID:id,Name:name} }
func (c *CatalogNode8_100) Rename(n string) { c.mu.Lock(); defer c.mu.Unlock(); c.Name=n }
func (c *CatalogNode8_100) String() string { return fmt.Sprintf("CatalogNode8_100 %s %s", c.ID, c.Name) }


