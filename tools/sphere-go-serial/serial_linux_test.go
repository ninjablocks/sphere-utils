package main

import (
	"testing"
)

func TestKernelCmdlineExtract(t *testing.T) {
	_, err := extractSerialFromCmdline([]byte("testing without hw serial = blah"))
	if err == nil {
		t.Fatal("Expected error while parsing invalid cmdline")
	}

	_, err = extractSerialFromCmdline([]byte("testing without hwserial=0101"))
	if err == nil {
		t.Fatal("Expected error while parsing cmdline with short serial")
	}

	_, err = extractSerialFromCmdline([]byte("testing without hwserial=testing01010101010101010101010101010101"))
	if err == nil {
		t.Fatal("Expected error while parsing cmdline with invalid hex characters in serial")
	}

	ser, err := extractSerialFromCmdline([]byte("console=ttyO0,115200n8 hwserial=aabbccddeeff9988ff00ff00ff00ff00 root=/dev/disk/by-label/system-a init=/lib/systemd/systemd ro panic=-1 fixrtc rootfstype=ext4"))
	if err != nil {
		t.Fatal("Expected successful serial from cmdline")
	}

	if ser != "VK54ZXPO76MYQ" {
		t.Fatalf("Serial %v does not match expected %v", ser, "VK54ZXPO76MYQ")
	}
}
