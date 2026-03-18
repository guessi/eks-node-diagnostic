package validate

import (
	"testing"
)

func TestNodeName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// instance ID: valid
		{"bare instance id", "i-0123456789abcdef0", false},

		// instance ID: invalid
		{"instance id too short", "i-0123456789abcdef", true},
		{"instance id too long", "i-0123456789abcdef01", true},
		{"instance id uppercase hex", "i-0123456789ABCDEF0", true},
		{"instance id wrong prefix", "x-0123456789abcdef0", true},
		{"instance id with ec2.internal", "i-0123456789abcdef0.ec2.internal", true},
		{"instance id with compute.internal", "i-0123456789abcdef0.us-west-2.compute.internal", true},

		// ip-name: valid
		{"ip name with ec2.internal", "ip-10-0-0-1.ec2.internal", false},
		{"ip name with us-west-2", "ip-10-24-34-0.us-west-2.compute.internal", false},
		{"ip name with us-east-2", "ip-10-0-0-1.us-east-2.compute.internal", false},
		{"ip name with ap-southeast-1", "ip-172-31-0-5.ap-southeast-1.compute.internal", false},
		{"ip name with max octets", "ip-255-255-255-255.ec2.internal", false},
		{"ip name with min first octet", "ip-1-0-0-0.ec2.internal", false},

		// ip-name: invalid IPv4
		{"ip name with octet over 255", "ip-999-0-0-1.ec2.internal", true},
		{"ip name with leading zero", "ip-010-0-0-1.ec2.internal", true},
		{"ip name with first octet zero", "ip-0-0-0-1.ec2.internal", true},

		// ip-name: invalid suffix
		{"ip name with us-east-1 compute suffix", "ip-10-0-0-1.us-east-1.compute.internal", true},
		{"ip name with wrong suffix", "ip-10-0-0-1.example.com", true},
		{"ip name missing .internal", "ip-10-0-0-1.ec2", true},
		{"ip name with fake region", "ip-10-0-0-1.xx-east-1.compute.internal", true},
		{"ip name with region number zero", "ip-10-0-0-1.ap-east-0.compute.internal", true},

		// other
		{"empty input", "", true},
		{"invalid node name", "not-a-valid-node", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NodeName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
