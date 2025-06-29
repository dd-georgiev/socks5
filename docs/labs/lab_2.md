- [Lab 2 - Implementing messages](#lab-2---implementing-messages)
    * [Overview](#overview)
    * [Concepts](#concepts)
        + [Encapsulation](#encapsulation)
        + [Encoding](#encoding)
        + [Message implementation standard overview](#message-implementation-standard-overview)
    * [Example - Implementing the Message interface for "available authentication methods" message](#example---implementing-the-message-interface-for--available-authentication-methods--message)
        + [Implementing the ToBytes function](#implementing-the-tobytes-function)
        + [Implementing the Deserialize function](#implementing-the-deserialize-function)
        + [Unit testing good paths](#unit-testing-good-paths)
        + [Benchmarking](#benchmarking)
        + [Fuzzing](#fuzzing)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>

# Lab 2 - Implementing messages
## Overview
The idea here it to implement the messages from lab 1. In Go. The messages can be roughly divided into three categories - `requests`, `response` and `data encapsulation`. Each message must be serializable/deserializable. The goal is being able to create all the messages in the lab. To accomplish this:

Create a method capable of deserializing data received via network to the structure
Each individual structure must be serializable with `ToBytes` method

## Concepts
### Encapsulation
Encapsulation is common practice, it can be observed every time a higher-level protocol uses the services of lower-level. For example TCP segments are encapsulated in packages by the IP protocol. The socks5 proxy encapsulates the data between the client and the proxy server, but the data to the end-server is not subject to encapsulation.

In SOCKS5, encapsulation is used when the data is being transmitted between the client and the proxy. In lab 1, as well as in the RFC it can be observed that the UDP traffic must be enclosed in specific message, where at the beginning of it, there is information about how it should be proxies, and the data its self is appended at the end.
### Encoding
The data transmitted over a network rarely matches the data in memory, because of this a set of encoding practices are commonly uses. For example the HTTP protocol uses `JSON` for communication between backend and front-end. The Postgres protocol uses a TLV(Type-Length Value) to send receive commands and transfer data. Certificates are encoded using X509. There are many other ways to encode data. The core idea in encoding is to encode the in-memory data structure in binary in such way, so that it can be reconstructed at the other end correctly. SOCKS5 uses fairly static structures for encoding information, most of the fields are with predefined length. The two exceptions are when the `DST.ADDR` field is `FQDN` and the data during encapsulation.

### Message implementation standard overview
In order for the implementation to be successful:

1. It must implement the Message interface, by implementing the `ToBytes` and `Deserialize` methods.
2. It must be tested that it returns error if the VER field is incorrect, or the auth methods are invalid(e.g. the number is not assigned by IANA)
3. It must have tests for benchmarking the ToBytes and Deserialize methods.
4. It must be fuzz-tested, interesting scenarios must be added to the fuzzing function, if any. Fuzzing for a few seconds(30) should be enough.
5. The “happy paths” must be covered by unit tests
6. Each file must contain information about the message, at the very least the message format, quoted from the results from Lab-1

The message interface is defined as follows:
```go
type Socks5Message interface {
	ToByte() ([]byte, error)
	Deserialize([]byte) error
}
```

In the next section, the implementation for the first message exchanged by the protocol (available authentication methods) will be implemented
## Example - Implementing the Message interface for "available authentication methods" message

In the next two sections, I will cover the implementation of the `Message` interface, which states the each message must have two methods - `ToBytes` and `Deserialize`. The idea of those methods is to either transform from or to binary/bytes. This conversation is done without any new instances of the struct. That is, the `Deserialize` method is mutating the structure. The `ToBytes` is returning a byte array, without modifying the instance in any way.

The struct looks like this:

```go
type AvailableAuthMethods struct {
	methods []uint16 // note that this field is private
}

// Getter
func (m *AvailableAuthMethods) Methods() []uint16 {
	return m.methods
}
// Setter, basically checks if the message is known(i.e. defined by IANA), if so it appends it, error otherwise
func (m *AvailableAuthMethods) AddMethods(method uint16) error {
	if method > JsonParameterBlock || method == Unassigned {
		return messages.UnknownAuthMethodError{Method: method}
	}
	m.methods = append(m.methods, method)
	return nil
}
```

### Implementing the ToBytes function
```go
func (m *AvailableAuthMethods) ToBytes() []byte {
	typesBytes := make([]byte, 0)
	for _, method := range m.methods {
		typesBytes = append(typesBytes, byte(method))
	}
	headersBytes := []byte{messages.PROTOCOL_VERSION, byte(len(m.methods))}
	return append(headersBytes, typesBytes...)
}
```
### Implementing the Deserialize function
```go
func (m *AvailableAuthMethods) Deserialize(buf []byte) error {
	if len(buf) < 3 { // it doesn't make sense to have less than that, as a single auth method + length + protocol version will be 3 bytes
		return messages.MalformedMessageError{}
	}
	if buf[MESSAGE_VERSION_INDEX] != messages.PROTOCOL_VERSION {
		return messages.MismatchedSocksVersionError{}
	}
	authMethodsCount := uint16(buf[MESSAGE_AVAIL_METHODS_INDEX])
	lastAuthMethodIndex := int(MESSAGE_AUTH_METHODS_START_INDEX + authMethodsCount)
	if lastAuthMethodIndex >= len(buf) { // make sure we are not out reaching the array
		return messages.MalformedMessageError{}
	}
	for i := MESSAGE_AUTH_METHODS_START_INDEX; i < lastAuthMethodIndex; i++ {
		currentAuthMethod := uint16(buf[i])
		if err := m.AddMethods(currentAuthMethod); err != nil {
			return err
		}
	}
	return nil
}
```

### Unit testing good paths
Only single unit test is covered, together with the primary helper function.

```go
func getCorrectBytes(methods []uint16) []byte {
	methodsByte := make([]byte, 0)
	for _, method := range methods {
		methodsByte = append(methodsByte, byte(method))
	}
	return append([]byte{0x05, byte(len(methods))}, methodsByte...)
}
func TestAvailableAuthMethods_Deserialize_Single_Method(t *testing.T) {
	validMethods := []uint16{0, 1, 2, 3, 5, 6, 7, 8, 9} // all except 4, as defined by IANA
	for _, method := range validMethods {
		singleMethod := AvailableAuthMethods{}
		err := singleMethod.Deserialize(getCorrectBytes([]uint16{method}))
		if err != nil {
			t.Fatal("Failed to deserialize correctly", method, err)
		}
	}
}
```
### Benchmarking
```go
func BenchmarkAvailableAuthMethods_Deserialize_Single_Method(b *testing.B) {
	req := []byte{0x05, 0x01, 0x01}
	for i := 0; i < b.N; i++ {
		msg := AvailableAuthMethods{}
		_ = msg.Deserialize(req)
	}
}
```
### Fuzzing
```go
func FuzzAvailableAuthMethods_Deserialize(f *testing.F) {
	f.Add([]byte{}) // this is how one adds specific condition
	f.Add([]byte{0x00, 0x05, 0x01}) // in case the size of the methods doesn't match the actual size send
	f.Fuzz(func(t *testing.T, data []byte) {
		msg := AvailableAuthMethods{}
		err := msg.Deserialize(data)
		if err != nil && !isKnownError(err) {
			t.Fatalf("Unexpected error %v with data %+v", err, data)
		}
	})
}
// quite ugly function which checks if the error is expected. Those were created by manually reasoning about what must be returned.
// For more complex functions it can be done based on experimentation(run fuzzing, check if error is ok, add to list or fix)
func isKnownError(err error) bool {
	return strings.Contains(err.Error(), "Mismatched socks version") ||
		strings.Contains(err.Error(), "Unknown auth method") ||
		strings.Contains(err.Error(), "Message is malformed")
}
```