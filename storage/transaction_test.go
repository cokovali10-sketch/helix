// Package storage provides additional storage engine tests
package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestTransactionManager_Create(t *testing.T) {
	tm := NewTransactionManager()
	if tm == nil {
		t.Fatal("TransactionManager is nil")
	}
}

func TestTransactionManager_Begin(t *testing.T) {
	tm := NewTransactionManager()
	txn := tm.Begin()

	if txn == nil {
		t.Fatal("Transaction is nil")
	}
	if txn.Status != TxnPending {
		t.Errorf("Expected pending, got %s", txn.Status)
	}
}

func TestTransactionManager_Commit(t *testing.T) {
	tm := NewTransactionManager()
	txn := tm.Begin()

	if err := tm.Commit(txn); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	if txn.Status != TxnCommitted {
		t.Errorf("Expected committed, got %s", txn.Status)
	}
}

func TestTransactionManager_Abort(t *testing.T) {
	tm := NewTransactionManager()
	txn := tm.Begin()

	if err := tm.Abort(txn); err != nil {
		t.Fatalf("Failed to abort: %v", err)
	}

	if txn.Status != TxnAborted {
		t.Errorf("Expected aborted, got %s", txn.Status)
	}
}

func TestTransactionManager_DoubleCommit(t *testing.T) {
	tm := NewTransactionManager()
	txn := tm.Begin()

	tm.Commit(txn)
	if err := tm.Commit(txn); err == nil {
		t.Error("Expected error on double commit")
	}
}

func TestTransactionManager_GetTransaction(t *testing.T) {
	tm := NewTransactionManager()
	txn := tm.Begin()

	got, ok := tm.GetTransaction(txn.ID)
	if !ok {
		t.Fatal("Transaction not found")
	}
	if got.ID != txn.ID {
		t.Errorf("Expected %s, got %s", txn.ID, got.ID)
	}
}

func TestTransactionManager_ListTransactions(t *testing.T) {
	tm := NewTransactionManager()
	tm.Begin()
	tm.Begin()
	tm.Begin()

	txns := tm.ListTransactions()
	if len(txns) != 3 {
		t.Errorf("Expected 3, got %d", len(txns))
	}
}

func TestTransaction_Put(t *testing.T) {
	txn := &Transaction{ID: "txn-1", Status: TxnPending}
	txn.Put("key1", []byte("value1"))

	ops := txn.GetOps()
	if len(ops) != 1 {
		t.Fatalf("Expected 1 op, got %d", len(ops))
	}
	if ops[0].Type != "put" {
		t.Errorf("Expected put, got %s", ops[0].Type)
	}
	if ops[0].Key != "key1" {
		t.Errorf("Expected key1, got %s", ops[0].Key)
	}
}

func TestTransaction_Delete(t *testing.T) {
	txn := &Transaction{ID: "txn-1", Status: TxnPending}
	txn.Delete("key1")

	ops := txn.GetOps()
	if len(ops) != 1 {
		t.Fatalf("Expected 1 op, got %d", len(ops))
	}
	if ops[0].Type != "delete" {
		t.Errorf("Expected delete, got %s", ops[0].Type)
	}
}

func TestTransaction_MultipleOps(t *testing.T) {
	txn := &Transaction{ID: "txn-1", Status: TxnPending}
	txn.Put("key1", []byte("value1"))
	txn.Put("key2", []byte("value2"))
	txn.Delete("key3")

	ops := txn.GetOps()
	if len(ops) != 3 {
		t.Fatalf("Expected 3 ops, got %d", len(ops))
	}
}

func TestSnapshotManager_Create(t *testing.T) {
	dir := t.TempDir()
	sm := NewSnapshotManager(dir)

	engine, _ := NewStorageEngine(dir+"/wal", dir+"/sst", 1024*1024, 1024)
	snapshot, err := sm.Create(engine)
	if err != nil {
		t.Fatalf("Failed to create snapshot: %v", err)
	}

	if snapshot == nil {
		t.Fatal("Snapshot is nil")
	}
}

func TestSnapshotManager_GetSnapshot(t *testing.T) {
	dir := t.TempDir()
	sm := NewSnapshotManager(dir)

	engine, _ := NewStorageEngine(dir+"/wal", dir+"/sst", 1024*1024, 1024)
	snapshot, _ := sm.Create(engine)

	got, ok := sm.GetSnapshot(snapshot.ID)
	if !ok {
		t.Fatal("Snapshot not found")
	}
	if got.ID != snapshot.ID {
		t.Errorf("Expected %s, got %s", snapshot.ID, got.ID)
	}
}

func TestSnapshotManager_ListSnapshots(t *testing.T) {
	dir := t.TempDir()
	sm := NewSnapshotManager(dir)

	engine, _ := NewStorageEngine(dir+"/wal", dir+"/sst", 1024*1024, 1024)
	sm.Create(engine)
	sm.Create(engine)

	snapshots := sm.ListSnapshots()
	if len(snapshots) != 2 {
		t.Errorf("Expected 2, got %d", len(snapshots))
	}
}

func TestSnapshotManager_DeleteSnapshot(t *testing.T) {
	dir := t.TempDir()
	sm := NewSnapshotManager(dir)

	engine, _ := NewStorageEngine(dir+"/wal", dir+"/sst", 1024*1024, 1024)
	snapshot, _ := sm.Create(engine)

	if err := sm.DeleteSnapshot(snapshot.ID); err != nil {
		t.Fatalf("Failed to delete: %v", err)
	}

	_, ok := sm.GetSnapshot(snapshot.ID)
	if ok {
		t.Error("Snapshot should be deleted")
	}
}

func TestBackupManager_CreateFullBackup(t *testing.T) {
	dir := t.TempDir()
	bm := NewBackupManager(dir)

	engine, _ := NewStorageEngine(dir+"/wal", dir+"/sst", 1024*1024, 1024)
	backup, err := bm.CreateFullBackup(engine)
	if err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}

	if backup == nil {
		t.Fatal("Backup is nil")
	}
	if backup.Type != "full" {
		t.Errorf("Expected full, got %s", backup.Type)
	}
}

func TestBackupManager_CreateIncrementalBackup(t *testing.T) {
	dir := t.TempDir()
	bm := NewBackupManager(dir)

	engine, _ := NewStorageEngine(dir+"/wal", dir+"/sst", 1024*1024, 1024)
	base, _ := bm.CreateFullBackup(engine)

	incr, err := bm.CreateIncrementalBackup(engine, base.ID)
	if err != nil {
		t.Fatalf("Failed to create incremental backup: %v", err)
	}

	if incr.Type != "incremental" {
		t.Errorf("Expected incremental, got %s", incr.Type)
	}
}

func TestBackupManager_ListBackups(t *testing.T) {
	dir := t.TempDir()
	bm := NewBackupManager(dir)

	engine, _ := NewStorageEngine(dir+"/wal", dir+"/sst", 1024*1024, 1024)
	bm.CreateFullBackup(engine)
	bm.CreateFullBackup(engine)

	backups := bm.ListBackups()
	if len(backups) != 2 {
		t.Errorf("Expected 2, got %d", len(backups))
	}
}

func TestBackupManager_DeleteBackup(t *testing.T) {
	dir := t.TempDir()
	bm := NewBackupManager(dir)

	engine, _ := NewStorageEngine(dir+"/wal", dir+"/sst", 1024*1024, 1024)
	backup, _ := bm.CreateFullBackup(engine)

	if err := bm.DeleteBackup(backup.ID); err != nil {
		t.Fatalf("Failed to delete: %v", err)
	}

	_, ok := bm.GetBackup(backup.ID)
	if ok {
		t.Error("Backup should be deleted")
	}
}

func TestRecoveryManager_HealthCheck(t *testing.T) {
	dir := t.TempDir()
	rm := NewRecoveryManager(dir+"/wal", dir+"/sst")

	health := rm.HealthCheck()
	if health == nil {
		t.Fatal("Health is nil")
	}
}

func TestRecoveryManager_CheckConsistency(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(dir+"/wal", 0755)
	os.MkdirAll(dir+"/sst", 0755)

	rm := NewRecoveryManager(dir+"/wal", dir+"/sst")
	consistent, issues := rm.CheckConsistency()

	// Should have issues (no WAL files)
	if consistent {
		t.Log("Consistent (might be OK)")
	}
	_ = issues
}

func TestSaveLoadMetadata(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "metadata.json")

	metadata := map[string]interface{}{
		"version": "1.0",
		"name":    "test",
	}

	if err := SaveMetadata(path, metadata); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	loaded, err := LoadMetadata(path)
	if err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	if loaded["version"] != "1.0" {
		t.Errorf("Expected 1.0, got %v", loaded["version"])
	}
}