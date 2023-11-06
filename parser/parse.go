package parser

import (
	"bytes"
	"encoding/base64"
	"net/url"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/doranych/gaexporter/protos"
)

func GetMigrationPayload(txt string) (*protos.MigrationPayload, error) {
	u, err := url.Parse(txt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse input")
	}
	data := []byte(u.Query().Get("data"))
	b := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	_, err = base64.StdEncoding.Decode(b, data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode input")
	}
	b = bytes.Trim(b, "\x00")

	payload := new(protos.MigrationPayload)
	if err = proto.Unmarshal(b, payload); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal input")
	}
	return payload, nil
}
