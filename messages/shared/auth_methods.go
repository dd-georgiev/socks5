package shared

// A name->uint16 mapping for the authentication methods defined in IANA.
//
// Their presence doesn't guarantee the availability or implementation on the server. As such the absence of a given method must be handled by higher-level abstraction.
//
// Reference: https://www.iana.org/assignments/socks-methods/socks-methods.xhtml
const (
	NoAuthRequired      = 0
	GSSAPI              = 1
	UsernameAndPassword = 2
	CHAP                = 3
	Unassigned          = 4
	CRAM                = 5
	SSL                 = 6
	NDS                 = 7
	MAF                 = 8
	JsonParameterBlock  = 9
	NoAcceptableMethods = 255
)
