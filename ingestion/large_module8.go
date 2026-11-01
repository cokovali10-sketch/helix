// Package ingestion provides large module 8
package ingestion

import (
	"fmt"
	"sync"
)

// IngestNode8_1
type IngestNode8_1 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_1(id string) *IngestNode8_1 { return &IngestNode8_1{ID:id} }
func (n *IngestNode8_1) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_1) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_1) String() string { return fmt.Sprintf("IngestNode8_1 %s %d", n.ID, n.Rate) }

// IngestNode8_2
type IngestNode8_2 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_2(id string) *IngestNode8_2 { return &IngestNode8_2{ID:id} }
func (n *IngestNode8_2) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_2) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_2) String() string { return fmt.Sprintf("IngestNode8_2 %s %d", n.ID, n.Rate) }

// IngestNode8_3
type IngestNode8_3 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_3(id string) *IngestNode8_3 { return &IngestNode8_3{ID:id} }
func (n *IngestNode8_3) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_3) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_3) String() string { return fmt.Sprintf("IngestNode8_3 %s %d", n.ID, n.Rate) }

// IngestNode8_4
type IngestNode8_4 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_4(id string) *IngestNode8_4 { return &IngestNode8_4{ID:id} }
func (n *IngestNode8_4) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_4) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_4) String() string { return fmt.Sprintf("IngestNode8_4 %s %d", n.ID, n.Rate) }

// IngestNode8_5
type IngestNode8_5 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_5(id string) *IngestNode8_5 { return &IngestNode8_5{ID:id} }
func (n *IngestNode8_5) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_5) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_5) String() string { return fmt.Sprintf("IngestNode8_5 %s %d", n.ID, n.Rate) }

// IngestNode8_6
type IngestNode8_6 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_6(id string) *IngestNode8_6 { return &IngestNode8_6{ID:id} }
func (n *IngestNode8_6) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_6) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_6) String() string { return fmt.Sprintf("IngestNode8_6 %s %d", n.ID, n.Rate) }

// IngestNode8_7
type IngestNode8_7 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_7(id string) *IngestNode8_7 { return &IngestNode8_7{ID:id} }
func (n *IngestNode8_7) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_7) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_7) String() string { return fmt.Sprintf("IngestNode8_7 %s %d", n.ID, n.Rate) }

// IngestNode8_8
type IngestNode8_8 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_8(id string) *IngestNode8_8 { return &IngestNode8_8{ID:id} }
func (n *IngestNode8_8) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_8) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_8) String() string { return fmt.Sprintf("IngestNode8_8 %s %d", n.ID, n.Rate) }

// IngestNode8_9
type IngestNode8_9 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_9(id string) *IngestNode8_9 { return &IngestNode8_9{ID:id} }
func (n *IngestNode8_9) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_9) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_9) String() string { return fmt.Sprintf("IngestNode8_9 %s %d", n.ID, n.Rate) }

// IngestNode8_10
type IngestNode8_10 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_10(id string) *IngestNode8_10 { return &IngestNode8_10{ID:id} }
func (n *IngestNode8_10) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_10) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_10) String() string { return fmt.Sprintf("IngestNode8_10 %s %d", n.ID, n.Rate) }

// IngestNode8_11
type IngestNode8_11 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_11(id string) *IngestNode8_11 { return &IngestNode8_11{ID:id} }
func (n *IngestNode8_11) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_11) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_11) String() string { return fmt.Sprintf("IngestNode8_11 %s %d", n.ID, n.Rate) }

// IngestNode8_12
type IngestNode8_12 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_12(id string) *IngestNode8_12 { return &IngestNode8_12{ID:id} }
func (n *IngestNode8_12) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_12) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_12) String() string { return fmt.Sprintf("IngestNode8_12 %s %d", n.ID, n.Rate) }

// IngestNode8_13
type IngestNode8_13 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_13(id string) *IngestNode8_13 { return &IngestNode8_13{ID:id} }
func (n *IngestNode8_13) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_13) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_13) String() string { return fmt.Sprintf("IngestNode8_13 %s %d", n.ID, n.Rate) }

// IngestNode8_14
type IngestNode8_14 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_14(id string) *IngestNode8_14 { return &IngestNode8_14{ID:id} }
func (n *IngestNode8_14) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_14) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_14) String() string { return fmt.Sprintf("IngestNode8_14 %s %d", n.ID, n.Rate) }

// IngestNode8_15
type IngestNode8_15 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_15(id string) *IngestNode8_15 { return &IngestNode8_15{ID:id} }
func (n *IngestNode8_15) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_15) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_15) String() string { return fmt.Sprintf("IngestNode8_15 %s %d", n.ID, n.Rate) }

// IngestNode8_16
type IngestNode8_16 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_16(id string) *IngestNode8_16 { return &IngestNode8_16{ID:id} }
func (n *IngestNode8_16) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_16) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_16) String() string { return fmt.Sprintf("IngestNode8_16 %s %d", n.ID, n.Rate) }

// IngestNode8_17
type IngestNode8_17 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_17(id string) *IngestNode8_17 { return &IngestNode8_17{ID:id} }
func (n *IngestNode8_17) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_17) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_17) String() string { return fmt.Sprintf("IngestNode8_17 %s %d", n.ID, n.Rate) }

// IngestNode8_18
type IngestNode8_18 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_18(id string) *IngestNode8_18 { return &IngestNode8_18{ID:id} }
func (n *IngestNode8_18) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_18) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_18) String() string { return fmt.Sprintf("IngestNode8_18 %s %d", n.ID, n.Rate) }

// IngestNode8_19
type IngestNode8_19 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_19(id string) *IngestNode8_19 { return &IngestNode8_19{ID:id} }
func (n *IngestNode8_19) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_19) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_19) String() string { return fmt.Sprintf("IngestNode8_19 %s %d", n.ID, n.Rate) }

// IngestNode8_20
type IngestNode8_20 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_20(id string) *IngestNode8_20 { return &IngestNode8_20{ID:id} }
func (n *IngestNode8_20) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_20) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_20) String() string { return fmt.Sprintf("IngestNode8_20 %s %d", n.ID, n.Rate) }

// IngestNode8_21
type IngestNode8_21 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_21(id string) *IngestNode8_21 { return &IngestNode8_21{ID:id} }
func (n *IngestNode8_21) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_21) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_21) String() string { return fmt.Sprintf("IngestNode8_21 %s %d", n.ID, n.Rate) }

// IngestNode8_22
type IngestNode8_22 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_22(id string) *IngestNode8_22 { return &IngestNode8_22{ID:id} }
func (n *IngestNode8_22) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_22) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_22) String() string { return fmt.Sprintf("IngestNode8_22 %s %d", n.ID, n.Rate) }

// IngestNode8_23
type IngestNode8_23 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_23(id string) *IngestNode8_23 { return &IngestNode8_23{ID:id} }
func (n *IngestNode8_23) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_23) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_23) String() string { return fmt.Sprintf("IngestNode8_23 %s %d", n.ID, n.Rate) }

// IngestNode8_24
type IngestNode8_24 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_24(id string) *IngestNode8_24 { return &IngestNode8_24{ID:id} }
func (n *IngestNode8_24) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_24) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_24) String() string { return fmt.Sprintf("IngestNode8_24 %s %d", n.ID, n.Rate) }

// IngestNode8_25
type IngestNode8_25 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_25(id string) *IngestNode8_25 { return &IngestNode8_25{ID:id} }
func (n *IngestNode8_25) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_25) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_25) String() string { return fmt.Sprintf("IngestNode8_25 %s %d", n.ID, n.Rate) }

// IngestNode8_26
type IngestNode8_26 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_26(id string) *IngestNode8_26 { return &IngestNode8_26{ID:id} }
func (n *IngestNode8_26) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_26) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_26) String() string { return fmt.Sprintf("IngestNode8_26 %s %d", n.ID, n.Rate) }

// IngestNode8_27
type IngestNode8_27 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_27(id string) *IngestNode8_27 { return &IngestNode8_27{ID:id} }
func (n *IngestNode8_27) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_27) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_27) String() string { return fmt.Sprintf("IngestNode8_27 %s %d", n.ID, n.Rate) }

// IngestNode8_28
type IngestNode8_28 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_28(id string) *IngestNode8_28 { return &IngestNode8_28{ID:id} }
func (n *IngestNode8_28) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_28) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_28) String() string { return fmt.Sprintf("IngestNode8_28 %s %d", n.ID, n.Rate) }

// IngestNode8_29
type IngestNode8_29 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_29(id string) *IngestNode8_29 { return &IngestNode8_29{ID:id} }
func (n *IngestNode8_29) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_29) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_29) String() string { return fmt.Sprintf("IngestNode8_29 %s %d", n.ID, n.Rate) }

// IngestNode8_30
type IngestNode8_30 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_30(id string) *IngestNode8_30 { return &IngestNode8_30{ID:id} }
func (n *IngestNode8_30) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_30) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_30) String() string { return fmt.Sprintf("IngestNode8_30 %s %d", n.ID, n.Rate) }

// IngestNode8_31
type IngestNode8_31 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_31(id string) *IngestNode8_31 { return &IngestNode8_31{ID:id} }
func (n *IngestNode8_31) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_31) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_31) String() string { return fmt.Sprintf("IngestNode8_31 %s %d", n.ID, n.Rate) }

// IngestNode8_32
type IngestNode8_32 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_32(id string) *IngestNode8_32 { return &IngestNode8_32{ID:id} }
func (n *IngestNode8_32) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_32) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_32) String() string { return fmt.Sprintf("IngestNode8_32 %s %d", n.ID, n.Rate) }

// IngestNode8_33
type IngestNode8_33 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_33(id string) *IngestNode8_33 { return &IngestNode8_33{ID:id} }
func (n *IngestNode8_33) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_33) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_33) String() string { return fmt.Sprintf("IngestNode8_33 %s %d", n.ID, n.Rate) }

// IngestNode8_34
type IngestNode8_34 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_34(id string) *IngestNode8_34 { return &IngestNode8_34{ID:id} }
func (n *IngestNode8_34) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_34) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_34) String() string { return fmt.Sprintf("IngestNode8_34 %s %d", n.ID, n.Rate) }

// IngestNode8_35
type IngestNode8_35 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_35(id string) *IngestNode8_35 { return &IngestNode8_35{ID:id} }
func (n *IngestNode8_35) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_35) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_35) String() string { return fmt.Sprintf("IngestNode8_35 %s %d", n.ID, n.Rate) }

// IngestNode8_36
type IngestNode8_36 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_36(id string) *IngestNode8_36 { return &IngestNode8_36{ID:id} }
func (n *IngestNode8_36) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_36) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_36) String() string { return fmt.Sprintf("IngestNode8_36 %s %d", n.ID, n.Rate) }

// IngestNode8_37
type IngestNode8_37 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_37(id string) *IngestNode8_37 { return &IngestNode8_37{ID:id} }
func (n *IngestNode8_37) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_37) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_37) String() string { return fmt.Sprintf("IngestNode8_37 %s %d", n.ID, n.Rate) }

// IngestNode8_38
type IngestNode8_38 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_38(id string) *IngestNode8_38 { return &IngestNode8_38{ID:id} }
func (n *IngestNode8_38) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_38) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_38) String() string { return fmt.Sprintf("IngestNode8_38 %s %d", n.ID, n.Rate) }

// IngestNode8_39
type IngestNode8_39 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_39(id string) *IngestNode8_39 { return &IngestNode8_39{ID:id} }
func (n *IngestNode8_39) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_39) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_39) String() string { return fmt.Sprintf("IngestNode8_39 %s %d", n.ID, n.Rate) }

// IngestNode8_40
type IngestNode8_40 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_40(id string) *IngestNode8_40 { return &IngestNode8_40{ID:id} }
func (n *IngestNode8_40) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_40) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_40) String() string { return fmt.Sprintf("IngestNode8_40 %s %d", n.ID, n.Rate) }

// IngestNode8_41
type IngestNode8_41 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_41(id string) *IngestNode8_41 { return &IngestNode8_41{ID:id} }
func (n *IngestNode8_41) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_41) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_41) String() string { return fmt.Sprintf("IngestNode8_41 %s %d", n.ID, n.Rate) }

// IngestNode8_42
type IngestNode8_42 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_42(id string) *IngestNode8_42 { return &IngestNode8_42{ID:id} }
func (n *IngestNode8_42) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_42) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_42) String() string { return fmt.Sprintf("IngestNode8_42 %s %d", n.ID, n.Rate) }

// IngestNode8_43
type IngestNode8_43 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_43(id string) *IngestNode8_43 { return &IngestNode8_43{ID:id} }
func (n *IngestNode8_43) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_43) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_43) String() string { return fmt.Sprintf("IngestNode8_43 %s %d", n.ID, n.Rate) }

// IngestNode8_44
type IngestNode8_44 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_44(id string) *IngestNode8_44 { return &IngestNode8_44{ID:id} }
func (n *IngestNode8_44) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_44) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_44) String() string { return fmt.Sprintf("IngestNode8_44 %s %d", n.ID, n.Rate) }

// IngestNode8_45
type IngestNode8_45 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_45(id string) *IngestNode8_45 { return &IngestNode8_45{ID:id} }
func (n *IngestNode8_45) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_45) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_45) String() string { return fmt.Sprintf("IngestNode8_45 %s %d", n.ID, n.Rate) }

// IngestNode8_46
type IngestNode8_46 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_46(id string) *IngestNode8_46 { return &IngestNode8_46{ID:id} }
func (n *IngestNode8_46) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_46) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_46) String() string { return fmt.Sprintf("IngestNode8_46 %s %d", n.ID, n.Rate) }

// IngestNode8_47
type IngestNode8_47 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_47(id string) *IngestNode8_47 { return &IngestNode8_47{ID:id} }
func (n *IngestNode8_47) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_47) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_47) String() string { return fmt.Sprintf("IngestNode8_47 %s %d", n.ID, n.Rate) }

// IngestNode8_48
type IngestNode8_48 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_48(id string) *IngestNode8_48 { return &IngestNode8_48{ID:id} }
func (n *IngestNode8_48) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_48) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_48) String() string { return fmt.Sprintf("IngestNode8_48 %s %d", n.ID, n.Rate) }

// IngestNode8_49
type IngestNode8_49 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_49(id string) *IngestNode8_49 { return &IngestNode8_49{ID:id} }
func (n *IngestNode8_49) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_49) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_49) String() string { return fmt.Sprintf("IngestNode8_49 %s %d", n.ID, n.Rate) }

// IngestNode8_50
type IngestNode8_50 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_50(id string) *IngestNode8_50 { return &IngestNode8_50{ID:id} }
func (n *IngestNode8_50) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_50) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_50) String() string { return fmt.Sprintf("IngestNode8_50 %s %d", n.ID, n.Rate) }

// IngestNode8_51
type IngestNode8_51 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_51(id string) *IngestNode8_51 { return &IngestNode8_51{ID:id} }
func (n *IngestNode8_51) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_51) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_51) String() string { return fmt.Sprintf("IngestNode8_51 %s %d", n.ID, n.Rate) }

// IngestNode8_52
type IngestNode8_52 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_52(id string) *IngestNode8_52 { return &IngestNode8_52{ID:id} }
func (n *IngestNode8_52) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_52) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_52) String() string { return fmt.Sprintf("IngestNode8_52 %s %d", n.ID, n.Rate) }

// IngestNode8_53
type IngestNode8_53 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_53(id string) *IngestNode8_53 { return &IngestNode8_53{ID:id} }
func (n *IngestNode8_53) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_53) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_53) String() string { return fmt.Sprintf("IngestNode8_53 %s %d", n.ID, n.Rate) }

// IngestNode8_54
type IngestNode8_54 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_54(id string) *IngestNode8_54 { return &IngestNode8_54{ID:id} }
func (n *IngestNode8_54) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_54) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_54) String() string { return fmt.Sprintf("IngestNode8_54 %s %d", n.ID, n.Rate) }

// IngestNode8_55
type IngestNode8_55 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_55(id string) *IngestNode8_55 { return &IngestNode8_55{ID:id} }
func (n *IngestNode8_55) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_55) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_55) String() string { return fmt.Sprintf("IngestNode8_55 %s %d", n.ID, n.Rate) }

// IngestNode8_56
type IngestNode8_56 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_56(id string) *IngestNode8_56 { return &IngestNode8_56{ID:id} }
func (n *IngestNode8_56) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_56) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_56) String() string { return fmt.Sprintf("IngestNode8_56 %s %d", n.ID, n.Rate) }

// IngestNode8_57
type IngestNode8_57 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_57(id string) *IngestNode8_57 { return &IngestNode8_57{ID:id} }
func (n *IngestNode8_57) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_57) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_57) String() string { return fmt.Sprintf("IngestNode8_57 %s %d", n.ID, n.Rate) }

// IngestNode8_58
type IngestNode8_58 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_58(id string) *IngestNode8_58 { return &IngestNode8_58{ID:id} }
func (n *IngestNode8_58) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_58) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_58) String() string { return fmt.Sprintf("IngestNode8_58 %s %d", n.ID, n.Rate) }

// IngestNode8_59
type IngestNode8_59 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_59(id string) *IngestNode8_59 { return &IngestNode8_59{ID:id} }
func (n *IngestNode8_59) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_59) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_59) String() string { return fmt.Sprintf("IngestNode8_59 %s %d", n.ID, n.Rate) }

// IngestNode8_60
type IngestNode8_60 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_60(id string) *IngestNode8_60 { return &IngestNode8_60{ID:id} }
func (n *IngestNode8_60) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_60) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_60) String() string { return fmt.Sprintf("IngestNode8_60 %s %d", n.ID, n.Rate) }

// IngestNode8_61
type IngestNode8_61 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_61(id string) *IngestNode8_61 { return &IngestNode8_61{ID:id} }
func (n *IngestNode8_61) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_61) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_61) String() string { return fmt.Sprintf("IngestNode8_61 %s %d", n.ID, n.Rate) }

// IngestNode8_62
type IngestNode8_62 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_62(id string) *IngestNode8_62 { return &IngestNode8_62{ID:id} }
func (n *IngestNode8_62) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_62) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_62) String() string { return fmt.Sprintf("IngestNode8_62 %s %d", n.ID, n.Rate) }

// IngestNode8_63
type IngestNode8_63 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_63(id string) *IngestNode8_63 { return &IngestNode8_63{ID:id} }
func (n *IngestNode8_63) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_63) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_63) String() string { return fmt.Sprintf("IngestNode8_63 %s %d", n.ID, n.Rate) }

// IngestNode8_64
type IngestNode8_64 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_64(id string) *IngestNode8_64 { return &IngestNode8_64{ID:id} }
func (n *IngestNode8_64) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_64) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_64) String() string { return fmt.Sprintf("IngestNode8_64 %s %d", n.ID, n.Rate) }

// IngestNode8_65
type IngestNode8_65 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_65(id string) *IngestNode8_65 { return &IngestNode8_65{ID:id} }
func (n *IngestNode8_65) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_65) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_65) String() string { return fmt.Sprintf("IngestNode8_65 %s %d", n.ID, n.Rate) }

// IngestNode8_66
type IngestNode8_66 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_66(id string) *IngestNode8_66 { return &IngestNode8_66{ID:id} }
func (n *IngestNode8_66) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_66) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_66) String() string { return fmt.Sprintf("IngestNode8_66 %s %d", n.ID, n.Rate) }

// IngestNode8_67
type IngestNode8_67 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_67(id string) *IngestNode8_67 { return &IngestNode8_67{ID:id} }
func (n *IngestNode8_67) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_67) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_67) String() string { return fmt.Sprintf("IngestNode8_67 %s %d", n.ID, n.Rate) }

// IngestNode8_68
type IngestNode8_68 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_68(id string) *IngestNode8_68 { return &IngestNode8_68{ID:id} }
func (n *IngestNode8_68) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_68) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_68) String() string { return fmt.Sprintf("IngestNode8_68 %s %d", n.ID, n.Rate) }

// IngestNode8_69
type IngestNode8_69 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_69(id string) *IngestNode8_69 { return &IngestNode8_69{ID:id} }
func (n *IngestNode8_69) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_69) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_69) String() string { return fmt.Sprintf("IngestNode8_69 %s %d", n.ID, n.Rate) }

// IngestNode8_70
type IngestNode8_70 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_70(id string) *IngestNode8_70 { return &IngestNode8_70{ID:id} }
func (n *IngestNode8_70) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_70) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_70) String() string { return fmt.Sprintf("IngestNode8_70 %s %d", n.ID, n.Rate) }

// IngestNode8_71
type IngestNode8_71 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_71(id string) *IngestNode8_71 { return &IngestNode8_71{ID:id} }
func (n *IngestNode8_71) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_71) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_71) String() string { return fmt.Sprintf("IngestNode8_71 %s %d", n.ID, n.Rate) }

// IngestNode8_72
type IngestNode8_72 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_72(id string) *IngestNode8_72 { return &IngestNode8_72{ID:id} }
func (n *IngestNode8_72) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_72) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_72) String() string { return fmt.Sprintf("IngestNode8_72 %s %d", n.ID, n.Rate) }

// IngestNode8_73
type IngestNode8_73 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_73(id string) *IngestNode8_73 { return &IngestNode8_73{ID:id} }
func (n *IngestNode8_73) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_73) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_73) String() string { return fmt.Sprintf("IngestNode8_73 %s %d", n.ID, n.Rate) }

// IngestNode8_74
type IngestNode8_74 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_74(id string) *IngestNode8_74 { return &IngestNode8_74{ID:id} }
func (n *IngestNode8_74) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_74) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_74) String() string { return fmt.Sprintf("IngestNode8_74 %s %d", n.ID, n.Rate) }

// IngestNode8_75
type IngestNode8_75 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_75(id string) *IngestNode8_75 { return &IngestNode8_75{ID:id} }
func (n *IngestNode8_75) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_75) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_75) String() string { return fmt.Sprintf("IngestNode8_75 %s %d", n.ID, n.Rate) }

// IngestNode8_76
type IngestNode8_76 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_76(id string) *IngestNode8_76 { return &IngestNode8_76{ID:id} }
func (n *IngestNode8_76) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_76) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_76) String() string { return fmt.Sprintf("IngestNode8_76 %s %d", n.ID, n.Rate) }

// IngestNode8_77
type IngestNode8_77 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_77(id string) *IngestNode8_77 { return &IngestNode8_77{ID:id} }
func (n *IngestNode8_77) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_77) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_77) String() string { return fmt.Sprintf("IngestNode8_77 %s %d", n.ID, n.Rate) }

// IngestNode8_78
type IngestNode8_78 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_78(id string) *IngestNode8_78 { return &IngestNode8_78{ID:id} }
func (n *IngestNode8_78) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_78) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_78) String() string { return fmt.Sprintf("IngestNode8_78 %s %d", n.ID, n.Rate) }

// IngestNode8_79
type IngestNode8_79 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_79(id string) *IngestNode8_79 { return &IngestNode8_79{ID:id} }
func (n *IngestNode8_79) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_79) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_79) String() string { return fmt.Sprintf("IngestNode8_79 %s %d", n.ID, n.Rate) }

// IngestNode8_80
type IngestNode8_80 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_80(id string) *IngestNode8_80 { return &IngestNode8_80{ID:id} }
func (n *IngestNode8_80) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_80) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_80) String() string { return fmt.Sprintf("IngestNode8_80 %s %d", n.ID, n.Rate) }

// IngestNode8_81
type IngestNode8_81 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_81(id string) *IngestNode8_81 { return &IngestNode8_81{ID:id} }
func (n *IngestNode8_81) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_81) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_81) String() string { return fmt.Sprintf("IngestNode8_81 %s %d", n.ID, n.Rate) }

// IngestNode8_82
type IngestNode8_82 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_82(id string) *IngestNode8_82 { return &IngestNode8_82{ID:id} }
func (n *IngestNode8_82) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_82) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_82) String() string { return fmt.Sprintf("IngestNode8_82 %s %d", n.ID, n.Rate) }

// IngestNode8_83
type IngestNode8_83 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_83(id string) *IngestNode8_83 { return &IngestNode8_83{ID:id} }
func (n *IngestNode8_83) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_83) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_83) String() string { return fmt.Sprintf("IngestNode8_83 %s %d", n.ID, n.Rate) }

// IngestNode8_84
type IngestNode8_84 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_84(id string) *IngestNode8_84 { return &IngestNode8_84{ID:id} }
func (n *IngestNode8_84) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_84) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_84) String() string { return fmt.Sprintf("IngestNode8_84 %s %d", n.ID, n.Rate) }

// IngestNode8_85
type IngestNode8_85 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_85(id string) *IngestNode8_85 { return &IngestNode8_85{ID:id} }
func (n *IngestNode8_85) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_85) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_85) String() string { return fmt.Sprintf("IngestNode8_85 %s %d", n.ID, n.Rate) }

// IngestNode8_86
type IngestNode8_86 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_86(id string) *IngestNode8_86 { return &IngestNode8_86{ID:id} }
func (n *IngestNode8_86) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_86) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_86) String() string { return fmt.Sprintf("IngestNode8_86 %s %d", n.ID, n.Rate) }

// IngestNode8_87
type IngestNode8_87 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_87(id string) *IngestNode8_87 { return &IngestNode8_87{ID:id} }
func (n *IngestNode8_87) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_87) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_87) String() string { return fmt.Sprintf("IngestNode8_87 %s %d", n.ID, n.Rate) }

// IngestNode8_88
type IngestNode8_88 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_88(id string) *IngestNode8_88 { return &IngestNode8_88{ID:id} }
func (n *IngestNode8_88) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_88) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_88) String() string { return fmt.Sprintf("IngestNode8_88 %s %d", n.ID, n.Rate) }

// IngestNode8_89
type IngestNode8_89 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_89(id string) *IngestNode8_89 { return &IngestNode8_89{ID:id} }
func (n *IngestNode8_89) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_89) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_89) String() string { return fmt.Sprintf("IngestNode8_89 %s %d", n.ID, n.Rate) }

// IngestNode8_90
type IngestNode8_90 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_90(id string) *IngestNode8_90 { return &IngestNode8_90{ID:id} }
func (n *IngestNode8_90) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_90) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_90) String() string { return fmt.Sprintf("IngestNode8_90 %s %d", n.ID, n.Rate) }

// IngestNode8_91
type IngestNode8_91 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_91(id string) *IngestNode8_91 { return &IngestNode8_91{ID:id} }
func (n *IngestNode8_91) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_91) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_91) String() string { return fmt.Sprintf("IngestNode8_91 %s %d", n.ID, n.Rate) }

// IngestNode8_92
type IngestNode8_92 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_92(id string) *IngestNode8_92 { return &IngestNode8_92{ID:id} }
func (n *IngestNode8_92) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_92) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_92) String() string { return fmt.Sprintf("IngestNode8_92 %s %d", n.ID, n.Rate) }

// IngestNode8_93
type IngestNode8_93 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_93(id string) *IngestNode8_93 { return &IngestNode8_93{ID:id} }
func (n *IngestNode8_93) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_93) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_93) String() string { return fmt.Sprintf("IngestNode8_93 %s %d", n.ID, n.Rate) }

// IngestNode8_94
type IngestNode8_94 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_94(id string) *IngestNode8_94 { return &IngestNode8_94{ID:id} }
func (n *IngestNode8_94) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_94) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_94) String() string { return fmt.Sprintf("IngestNode8_94 %s %d", n.ID, n.Rate) }

// IngestNode8_95
type IngestNode8_95 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_95(id string) *IngestNode8_95 { return &IngestNode8_95{ID:id} }
func (n *IngestNode8_95) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_95) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_95) String() string { return fmt.Sprintf("IngestNode8_95 %s %d", n.ID, n.Rate) }

// IngestNode8_96
type IngestNode8_96 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_96(id string) *IngestNode8_96 { return &IngestNode8_96{ID:id} }
func (n *IngestNode8_96) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_96) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_96) String() string { return fmt.Sprintf("IngestNode8_96 %s %d", n.ID, n.Rate) }

// IngestNode8_97
type IngestNode8_97 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_97(id string) *IngestNode8_97 { return &IngestNode8_97{ID:id} }
func (n *IngestNode8_97) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_97) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_97) String() string { return fmt.Sprintf("IngestNode8_97 %s %d", n.ID, n.Rate) }

// IngestNode8_98
type IngestNode8_98 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_98(id string) *IngestNode8_98 { return &IngestNode8_98{ID:id} }
func (n *IngestNode8_98) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_98) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_98) String() string { return fmt.Sprintf("IngestNode8_98 %s %d", n.ID, n.Rate) }

// IngestNode8_99
type IngestNode8_99 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_99(id string) *IngestNode8_99 { return &IngestNode8_99{ID:id} }
func (n *IngestNode8_99) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_99) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_99) String() string { return fmt.Sprintf("IngestNode8_99 %s %d", n.ID, n.Rate) }

// IngestNode8_100
type IngestNode8_100 struct {
	ID string
	Rate int
	mu sync.Mutex
}

func NewIngestNode8_100(id string) *IngestNode8_100 { return &IngestNode8_100{ID:id} }
func (n *IngestNode8_100) SetRate(r int) { n.mu.Lock(); defer n.mu.Unlock(); n.Rate=r }
func (n *IngestNode8_100) RateVal() int { n.mu.Lock(); defer n.mu.Unlock(); return n.Rate }
func (n *IngestNode8_100) String() string { return fmt.Sprintf("IngestNode8_100 %s %d", n.ID, n.Rate) }


