package logger

import (
	"strings"
	"time"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// customEncoder wraps a JSON encoder to format logs as specified.
type customEncoder struct {
	zapcore.Encoder
}

func (enc *customEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf := buffer.NewPool().Get()

	// Timestamp.
	buf.AppendString(entry.Time.Format("02/01/2006 15:04:05.000"))
	buf.AppendString(" | ")

	// Level.
	levelStr := strings.ToUpper(entry.Level.String())
	buf.AppendString(levelStr)
	buf.AppendString(" | ")

	// Caller.
	if entry.Caller.Defined {
		buf.AppendString(entry.Caller.TrimmedPath())
	} else {
		buf.AppendString("unknown")
	}
	buf.AppendString(" | ")

	// Trace if present.
	var traceID string
	var hasFields bool
	// remainingFields := make([]zapcore.Field, 0, len(fields))
	for _, f := range fields {
		if f.Key == "traceid" {
			if f.Type == zapcore.StringType {
				traceID = f.String
			}
			continue
		}
		if !hasFields {
			hasFields = true
		}
	}
	if traceID != "" {
		buf.AppendString("TRACE : ")
		buf.AppendString(traceID)
		buf.AppendString(" | ")
	}

	// Message.
	buf.AppendString(entry.Message)

	// Any other fields as key=val.
	if hasFields {
		for _, f := range fields {
			if f.Key == "traceid" {
				continue
			}
			buf.AppendByte(' ')
			buf.AppendString(f.Key)
			buf.AppendByte('=')
			f.AddTo(enc)
			// switch f.Type {
			// case zapcore.StringType:
			// 	buf.AppendByte('"')
			// 	buf.AppendString(f.String)
			// 	buf.AppendByte('"')
			// case zapcore.Int64Type, zapcore.Int32Type, zapcore.Uint64Type, zapcore.Uint32Type:
			// 	buf.AppendString(fmt.Sprint(f.Integer))
			// default:
			// 	buf.AppendString(fmt.Sprint(f.Interface))
			// }
		}
	}

	// Append the stacktrace if it exists.
	if entry.Stack != "" {
		// buf.AppendString("\n\x1b[35mSTACK\x1b[0m : \n")
		buf.AppendString("\nSTACK : \n")
		buf.AppendString(entry.Stack)
	}

	buf.AppendByte('\n')
	return buf, nil
}

func (enc *customEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	return enc.Encoder.AddArray(key, arr)
}
func (enc *customEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	return enc.Encoder.AddObject(key, obj)
}
func (enc *customEncoder) AddBinary(key string, val []byte) {
	enc.Encoder.AddBinary(key, val)
}
func (enc *customEncoder) AddByteString(key string, val []byte) {
	enc.Encoder.AddByteString(key, val)
}
func (enc *customEncoder) AddBool(key string, val bool) {
	enc.Encoder.AddBool(key, val)
}
func (enc *customEncoder) AddComplex128(key string, val complex128) {
	enc.Encoder.AddComplex128(key, val)
}
func (enc *customEncoder) AddComplex64(key string, val complex64) {
	enc.Encoder.AddComplex64(key, val)
}
func (enc *customEncoder) AddDuration(key string, val time.Duration) {
	enc.Encoder.AddDuration(key, val)
}
func (enc *customEncoder) AddFloat64(key string, val float64) {
	enc.Encoder.AddFloat64(key, val)
}
func (enc *customEncoder) AddFloat32(key string, val float32) {
	enc.Encoder.AddFloat32(key, val)
}
func (enc *customEncoder) AddInt64(key string, val int64) {
	enc.Encoder.AddInt64(key, val)
}
func (enc *customEncoder) AddInt32(key string, val int32) {
	enc.Encoder.AddInt32(key, val)
}
func (enc *customEncoder) AddInt16(key string, val int16) {
	enc.Encoder.AddInt16(key, val)
}
func (enc *customEncoder) AddInt8(key string, val int8) {
	enc.Encoder.AddInt8(key, val)
}
func (enc *customEncoder) AddString(key, val string) {
	enc.Encoder.AddString(key, val)
}
func (enc *customEncoder) AddTime(key string, val time.Time) {
	enc.Encoder.AddTime(key, val)
}
func (enc *customEncoder) AddUint64(key string, val uint64) {
	enc.Encoder.AddUint64(key, val)
}
func (enc *customEncoder) AddUint32(key string, val uint32) {
	enc.Encoder.AddUint32(key, val)
}
func (enc *customEncoder) AddUint16(key string, val uint16) {
	enc.Encoder.AddUint16(key, val)
}
func (enc *customEncoder) AddUint8(key string, val uint8) {
	enc.Encoder.AddUint8(key, val)
}
func (enc *customEncoder) AddUintptr(key string, val uintptr) {
	enc.Encoder.AddUintptr(key, val)
}
func (enc *customEncoder) AddReflected(key string, val any) error {
	return enc.Encoder.AddReflected(key, val)
}
func (enc *customEncoder) OpenNamespace(key string) {
	enc.Encoder.OpenNamespace(key)
}
func (enc *customEncoder) Clone() zapcore.Encoder {
	return enc.Encoder.Clone()
}
