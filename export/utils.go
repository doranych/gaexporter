package export

import (
	"encoding/base32"
	"net/url"
	"strconv"
	"strings"

	"github.com/doranych/gaexporter/parser"
	"github.com/pkg/errors"

	"github.com/doranych/gaexporter/protos"
)

func payloadToText(payload *protos.MigrationPayload) []byte {
	sb := strings.Builder{}
	for i, _ := range payload.GetOtpParameters() {
		sb.WriteString("Name: ")
		sb.WriteString(payload.GetOtpParameters()[i].GetName() + "\n")
		sb.WriteString("Secret: ")
		sb.WriteString(base32.StdEncoding.EncodeToString(payload.GetOtpParameters()[i].GetSecret()) + "\n")
		sb.WriteString("Type: ")
		sb.WriteString(otpType(int32(payload.GetOtpParameters()[i].GetType())) + "\n")
		sb.WriteString("Issuer: ")
		sb.WriteString(payload.GetOtpParameters()[i].GetIssuer() + "\n")
		sb.WriteString("Algorithm: ")
		sb.WriteString(otpAlgorithm(int32(payload.GetOtpParameters()[i].GetAlgorithm())) + "\n")
		sb.WriteString("Digits: ")
		sb.WriteString(strconv.Itoa(int(payload.GetOtpParameters()[i].GetDigits())) + "\n")
		sb.WriteString("Counter: ")
		sb.WriteString(strconv.Itoa(int(payload.GetOtpParameters()[i].GetCounter())) + "\n")
		sb.WriteString("OtpURL: ")
		sb.WriteString(otpUrl(payload.GetOtpParameters()[i]) + "\n")
		sb.WriteString("\n")
	}
	return []byte(sb.String())
}

func otpType(i int32) string {
	return strings.ToLower(strings.TrimPrefix(protos.MigrationPayload_OtpType_name[i], "OTP_TYPE_"))
}

func otpAlgorithm(i int32) string {
	return strings.ToLower(strings.TrimPrefix(protos.MigrationPayload_Algorithm_name[i], "ALGORITHM_"))
}

func otpUrl(payload *protos.MigrationPayload_OtpParameters) string {
	u := url.URL{
		Scheme: "otpauth",
		Host:   otpType(int32(payload.GetType())),
		Path:   url.PathEscape(payload.GetName()),
	}
	v := u.Query()
	if payload.GetType() == protos.MigrationPayload_OTP_TYPE_HOTP {
		v.Add("counter", strconv.Itoa(int(payload.GetCounter())))
	}
	if payload.GetIssuer() != "" {
		v.Add("issuer", payload.GetIssuer())
	}
	v.Add("secret", base32.StdEncoding.EncodeToString(payload.GetSecret()))
	u.RawQuery = v.Encode()
	return u.String()
}

func processMigrationUrl(str string, output Output, format string) error {
	payload, err := parser.GetMigrationPayload(str)
	if err != nil {
		return errors.Wrap(err, "failed to parse input")
	}

	switch format {
	case "txt":
		_, err = output.Dest.Write(payloadToText(payload))
		if err != nil {
			return errors.Wrap(err, "failed to write output")
		}
		return nil
	case "qr":
		_, err := output.Dest.Write(payloadToQr(payload))
		if err != nil {
			return errors.Wrap(err, "failed to write output")
		}
	default:
		return errors.New("format must be txt or qr")
	}
	return errors.New("not implemented")
}

func payloadToQr(payload *protos.MigrationPayload) []byte {
	return nil
}
