package auth

const (
	NO_AUTH                                     = 0
	GSSAPI                                      = 1
	USERNAME_PASSWORD                           = 2
	CHALLENGE_HANDSHAKE_AUTHENTICATION_PROTOCOL = 3
	UNASSIGNED                                  = 4
	CHALLENGE_RESPONSE_AUTHENTICATION_METHOD    = 5
	SECURE_SOCKETS_LAYER                        = 6
	NDS_AUTHENTICATION                          = 7
	MULTI_AUTHENTICATION_FRAMEWORK              = 8
	JSON_PARAMETER_BLOCK                        = 9
	NO_ACCEPTABLE_METHOD                        = 255
)

func IsValidAuthMethod(method uint16) bool {
	// since there are 9 methods registered in IANA and private methods are no supported, if the number is above 9 it references invalid method.
	return method < 10
}
