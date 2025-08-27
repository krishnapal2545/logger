package logger

import (
	"fmt"
	"strings"

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
	remainingFields := make([]zapcore.Field, 0, len(fields))
	for _, f := range fields {
		if f.Key == "traceid" && f.Type == zapcore.StringType {
			traceID = f.String
		} else {
			remainingFields = append(remainingFields, f)
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
	for _, f := range remainingFields {
		buf.AppendByte(' ')
		buf.AppendString(f.Key)
		buf.AppendByte('=')
		switch f.Type {
		case zapcore.StringType:
			buf.AppendByte('"')
			buf.AppendString(f.String)
			buf.AppendByte('"')
		case zapcore.Int64Type, zapcore.Int32Type, zapcore.Uint64Type, zapcore.Uint32Type:
			buf.AppendString(fmt.Sprint(f.Integer))
		default:
			buf.AppendString(fmt.Sprint(f.Interface))
		}
	}

	buf.AppendByte('\n')
	return buf, nil
}
