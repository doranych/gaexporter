package export

import (
	"bytes"
	"encoding/base32"
	"net/url"
	"strconv"
	"strings"

	"github.com/mdp/qrterminal/v3"
	"github.com/pkg/errors"
	"rsc.io/qr"

	"github.com/doranych/gaexporter/parser"
	"github.com/doranych/gaexporter/protos"
)

func payloadToText(payload *protos.MigrationPayload) []byte {
	sb := strings.Builder{}
	for i := range payload.GetOtpParameters() {
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
		qr, err := payloadToQr(payload, output)
		_, err = output.Dest.Write(qr)
		if err != nil {
			return errors.Wrap(err, "failed to write output")
		}
		return nil
	default:
		return errors.New("format must be txt or qr")
	}
}

func payloadToQr(payload *protos.MigrationPayload, output Output) ([]byte, error) {
	if output.Type == OutputTypeStdout {
		result := make([]byte, 0)

		result = append(result, payloadToText(payload)...)
		result = append(result, []byte("\n")...)

		for i, parameter := range payload.OtpParameters {
			buf := bytes.NewBuffer([]byte{})
			cfg := qrterminal.Config{
				Writer:    buf,
				QuietZone: 2,
				BlackChar: qrterminal.BLACK,
				WhiteChar: qrterminal.WHITE,
			}
			qrterminal.GenerateWithConfig(otpUrl(parameter), cfg)

			result = append(result, buf.Bytes()...)
			if i != len(payload.OtpParameters)-1 {
				result = append(result, []byte("\n")...)
			}
		}
		return result, nil
	}
	code, err := qr.Encode(otpUrl(payload.GetOtpParameters()[0]), qr.L)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode qr code")
	}
	return code.PNG(), err
}
