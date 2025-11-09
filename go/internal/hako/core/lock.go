// Package core provides locking mechanism for the Hakoniwa game.
// This file contains lock/unlock functions translated from Perl lib/Hako/Main.pm
//
// Ref: perl/lib/Hako/Main.pm:689-855
package core

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/neguse/hakoniwa/internal/hako/hconst"
	"github.com/neguse/hakoniwa/internal/hako/variable"
)

// hakoLock acquires a lock for data file access
// Ref: perl/lib/Hako/Main.pm:689
func hakoLock() bool {
	switch hconst.LockMode {
	case 1:
		return hakoLock1()
	case 2:
		return hakoLock2()
	case 3:
		return hakoLock3()
	default:
		return hakoLock4()
	}
}

// hakoLock1 uses directory-based locking
// Ref: perl/lib/Hako/Main.pm:712
func hakoLock1() bool {
	// Try to create lock directory
	err := os.Mkdir("hakojimalock", os.FileMode(hconst.DirMode))
	if err == nil {
		// Success
		return true
	}

	// Failed - check if should force unlock
	info, err := os.Stat("hakojimalock")
	if err == nil {
		modTime := info.ModTime()
		if time.Since(modTime).Seconds() > float64(hconst.UnlockTime) {
			// Force unlock
			unlock()
			// Phase 1: Should output tempUnlock() but simplified for now
			return false
		}
	}
	return false
}

// hakoLock2 uses flock-based locking (recommended)
// Ref: perl/lib/Hako/Main.pm:744
func hakoLock2() bool {
	// Open lock file
	file, err := os.OpenFile("hakojimalockflock", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return false
	}

	// Try to acquire exclusive lock
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		file.Close()
		return false
	}

	// Success - store file handle for later unlock
	variable.LockID = file
	return true
}

// hakoLock3 uses symlink-based locking
// Ref: perl/lib/Hako/Main.pm:757
func hakoLock3() bool {
	// Try to create symlink
	err := os.Symlink("hakojimalockdummy", "hakojimalock")
	if err == nil {
		// Success
		return true
	}

	// Failed - check if should force unlock
	info, err := os.Lstat("hakojimalock")
	if err == nil {
		modTime := info.ModTime()
		if time.Since(modTime).Seconds() > float64(hconst.UnlockTime) {
			// Force unlock
			unlock()
			// Phase 1: Should output tempUnlock() but simplified for now
			return false
		}
	}
	return false
}

// hakoLock4 uses file-based locking (not recommended)
// Ref: perl/lib/Hako/Main.pm:789
func hakoLock4() bool {
	// Try to remove key-free file
	err := os.Remove("key-free")
	if err == nil {
		// Success - create key-locked file
		file, err := os.Create("key-locked")
		if err == nil {
			fmt.Fprintf(file, "%d", time.Now().Unix())
			file.Close()
			return true
		}
	}

	// Failed - check if should force unlock
	data, err := os.ReadFile("key-locked")
	if err == nil {
		var lockTime int64
		fmt.Sscanf(string(data), "%d", &lockTime)
		if time.Now().Unix()-lockTime > int64(hconst.UnlockTime) {
			// Force unlock
			unlock()
			// Phase 1: Should output tempUnlock() but simplified for now
			return false
		}
	}
	return false
}

// Unlock releases the lock (exported for external use)
// Ref: perl/lib/Hako/Main.pm:832
func Unlock() {
	unlock()
}

// unlock releases the lock (internal)
func unlock() {
	switch hconst.LockMode {
	case 1:
		os.Remove("hakojimalock")
	case 2:
		if variable.LockID != nil {
			if file, ok := variable.LockID.(*os.File); ok {
				syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
				file.Close()
			}
			variable.LockID = nil
		}
	case 3:
		os.Remove("hakojimalock")
	default:
		os.Remove("key-locked")
		// Create key-free file
		file, err := os.Create("key-free")
		if err == nil {
			file.Close()
		}
	}
}
