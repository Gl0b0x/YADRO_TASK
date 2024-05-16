package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestComputerClub(t *testing.T) {
	var (
		expected1 = "09:00\n09:41 1 client1\n09:48 1 client2\n09:52 1 client3\n09:54 2 client1 1\n10:25 2 client2 2\n10:58 3 client3\n10:59 1 client4\n11:00 3 client4\n11:02 1 client5\n11:03 3 client5\n11:03 11 client5\n00:00 11 client1\n00:00 11 client2\n00:00 11 client3\n00:00 11 client4\n00:00 11 client5\n00:00\n1 150 14:06\n2 140 13:35\n"
		expected2 = "error: file not_file not found\n"
		expected3 = "error: command line arguments\n"
		expected4 = "09:00\n08:48 1 client1\n08:48 13 NotOpenYet\n09:41 1 client1\n09:48 1 client2\n09:52 3 client1\n09:52 13 ICanWaitNoLonger!\n09:54 2 client1 1\n10:25 2 client2 2\n10:58 1 client3\n10:59 2 client3 3\n11:30 1 client4\n11:35 2 client4 2\n11:35 13 PlaceIsBusy\n11:45 3 client4\n12:33 4 client1\n12:33 12 client4 1\n12:43 4 client2\n15:52 4 client4\n19:00 11 client3\n19:00\n1 70 05:58\n2 30 02:18\n3 90 08:01\n"
		expected5 = "error: not enough input data\n"
		expected6 = "-10\n"
		expected7 = "09:00 08:00\n"
		expected8 = "9:00 18:00\n"
	)
	testTable := []struct {
		testInfo string
		osArgs   []string
		expected string
	}{
		{"Overflow waiting List", []string{"", "test_file.txt"}, expected1},
		{"File not found", []string{"", "not_file"}, expected2},
		{"Incorrect command line arguments", []string{"", "test_file.txt", ""}, expected3},
		{"Example test", []string{"", "test_file1.txt"}, expected4},
		{"Test empty file", []string{"", "test_file2.txt"}, expected5},
		{"Incorrect countTables", []string{"", "test_file3.txt"}, expected6},
		{"Incorrect time work club", []string{"", "test_file4.txt"}, expected7},
		{"Invalid time work", []string{"", "test_file5.txt"}, expected8},
	}
testCase:
	for _, testCase := range testTable {
		t.Log(testCase.testInfo)
		os.Args = testCase.osArgs
		reader, writer, _ := os.Pipe()
		var buf bytes.Buffer
		os.Stdout = writer
		main()
		err := writer.Close()
		if err != nil {
			return
		}
		_, err = buf.ReadFrom(reader)
		if err != nil {
			return
		}
		output := strings.Split(buf.String(), "\n")
		expected := strings.Split(testCase.expected, "\n")
		if len(output) != len(expected) {
			t.Errorf("expected %d strings, got %d", len(expected), len(output))
			continue
		}
		for i := 0; i < len(output); i++ {
			if output[i] != expected[i] {
				t.Errorf("line %d expected %q, got %q", i+1, expected[i], output[i])
				continue testCase
			}
		}
	}
}
