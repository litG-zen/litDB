package parser

import "testing"

func TestSimpleString(t *testing.T) {
	data := []byte("+OK\r\n")
	value, err := Decode(data)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	if value != "OK" {
		t.Fatalf("Expected 'OK', got '%v'", value)
	}
}

func TestError(t *testing.T) {
	data := []byte("-ERR unknown command 'foo'\r\n")
	value, err := Decode(data)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	if value != "ERR unknown command 'foo'" {
		t.Fatalf("Expected 'ERR unknown command 'foo'', got '%v'", value)
	}
}

func TestInteger(t *testing.T) {
	data := []byte(":42\r\n")
	value, err := Decode(data)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	if value != int64(42) {
		t.Fatalf("Expected 42, got '%v'", value)
	}
}

func TestBulkString(t *testing.T) {
	data := []byte("$6\r\nfoobar\r\n")
	value, err := Decode(data)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	if value != "foobar" {
		t.Fatalf("Expected 'foobar', got '%v'", value)
	}
}

func TestArray(t *testing.T) {
	data := []byte("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n")
	value, err := Decode(data)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	arr, ok := value.([]interface{})
	if !ok || len(arr) != 2 {
		t.Fatalf("Expected array of length 2, got '%v'", value)
	}
	if arr[0] != "foo" || arr[1] != "bar" {
		t.Fatalf("Expected ['foo', 'bar'], got '%v'", arr)
	}
}
