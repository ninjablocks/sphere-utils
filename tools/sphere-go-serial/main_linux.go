package main

import (
	"bytes"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

func main() {

	flagset := flag.NewFlagSet("cmdline", flag.ContinueOnError)
	flagset.SetOutput(ioutil.Discard)
	hwSerial := flagset.String("hwserial", "", "The hardware serial")
	flagset.Parse(os.Args[1:])

	var err error
	input := []byte{}

	if hwSerial != nil && *hwSerial != "" {
		input = []byte(fmt.Sprintf("hwserial=%s ", *hwSerial))
	} else {
		input, err = ioutil.ReadFile("/proc/cmdline")
		if err != nil {
			panic(err)
		}
	}

	serial, err := extractSerialFromCmdline(input)
	if err != nil {
		serial, err = extractSerialFromBeagleEEPROM()
		if err != nil {
			serial, err = extractFromMacAddress()
			if err != nil {
				panic(err)
			}
		}
	}

	os.Stdout.Write([]byte(serial))
}

func extractSerialFromCmdline(cmdline []byte) (string, error) {
	parts := bytes.Split(cmdline, []byte("hwserial="))
	if len(parts) != 2 {
		return "", errors.New("Serial not present in kernel command line")
	}

	serialPart := parts[1]
	if len(serialPart) < 16 {
		return "", errors.New("Serial in kernel command line is too short")
	}

	rawSerial, err := hex.DecodeString(string(serialPart[0:16]))
	if err != nil {
		return "", err
	}

	b32Serial := strings.TrimRight(base32.StdEncoding.EncodeToString(rawSerial), "=")

	return b32Serial, nil
}

const i2cDeviceCreate = "/sys/bus/i2c/devices/i2c-0/new_device"
const i2CEEPROMDeviceNode = "/sys/bus/i2c/devices/0-0050/eeprom"

func extractSerialFromBeagleEEPROM() (string, error) {

	// this code needs to run setuid(root) to provide access to the necessary namespace

	if _, err := os.Stat(i2cDeviceCreate); os.IsNotExist(err) {
		return "", errors.New("I2C not available, likely not running on a BeagleBone")
	}

	// create EEPROM device if it doesn't exist
	if _, err := os.Stat(i2CEEPROMDeviceNode); os.IsNotExist(err) {
		toWrite := "24c256 0x50\n"

		err := ioutil.WriteFile(i2cDeviceCreate, []byte(toWrite), 0644)
		if err != nil {
			return "", err
		}
	}

	// xxd -g 2 -a -l 16 -seek 16 /sys/bus/i2c/devices/0-0050/eeprom | sed 's/^.* //' | sed -e 's/[.]//g' | tr -d '\n'
	file, err := os.Open(i2CEEPROMDeviceNode)
	if err != nil {
		return "", err
	}

	bytes := make([]byte, 16)
	_, err = file.ReadAt(bytes, 16)
	if err != nil {
		return "", err
	}

	file.Close()

	return strings.TrimRight(string(bytes), "\xff\x00"), nil
}

func extractFromMacAddress() (string, error) {

	// lock the goroutine to a known Linux thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// get the current state
	euid := os.Geteuid()
	uid := os.Getuid()
	egid := os.Getegid()
	gid := os.Getgid()

	// make the state safe
	if euid != uid {
		syscall.Setuid(uid)
	}
	if egid != gid {
		syscall.Setgid(gid)
	}

	// do the potentially unsafe stuff now...
	cmd := exec.Command("/bin/sh", "-c", "/sbin/ifconfig | /usr/bin/sed -n 's/\\([^ ]*\\).*HWaddr /\\1 /p' | /usr/bin/egrep '^(wlan|eth)' | /usr/bin/cut -f2 -d' ' | /usr/bin/tail -1 | /usr/bin/openssl md5 | /usr/bin/cut -f2 -d' ' | /usr/bin/base64 | /usr/bin/cut -c1-12 | /usr/bin/tr '[a-z]' 'A-Z'")
	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	} else {
		return "MAC" + string(bytes[0:len(bytes)-1]), nil
	}
}
