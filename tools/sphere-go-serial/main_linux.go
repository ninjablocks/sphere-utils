package main

import (
	"bytes"
	"encoding/base32"
	"encoding/hex"
	"io/ioutil"
	"os"
	"strings"
	"errors"
)

func main() {

	input, err := ioutil.ReadFile("/proc/cmdline")
	if err != nil {
		panic(err)
	}

	serial, err := extractSerialFromCmdline(input)
	if err != nil {
		serial, err = extractSerialFromBeagleEEPROM()
		if err != nil {
			panic(err)
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