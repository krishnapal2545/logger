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
	// Create a map for fields.
	fieldMap := make(map[string]any)
	for _, f := range fields {
		f.AddTo(zapcore.NewMapObjectEncoder())
	}

	buf := buffer.NewPool().Get()

	// Timestamp.
	buf.AppendString(entry.Time.Format("02-01-2006 15:04:05.000"))
	buf.AppendString(" | ")

	// Level.
	levelStr := strings.ToUpper(entry.Level.String())
	if entry.Level == zapcore.Level(-1) {
		levelStr = "TRACE"
	}
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
	if trace, ok := fieldMap["traceid"]; ok {
		buf.AppendString("TRACE : ")
		buf.AppendString(fmt.Sprint(trace))
		buf.AppendString(" | ")
		delete(fieldMap, "traceid")
	}

	// Message.
	buf.AppendString(entry.Message)

	// Any other fields as key=val.
	for k, v := range fieldMap {
		buf.AppendByte(' ')
		buf.AppendString(k)
		buf.AppendByte('=')
		if str, ok := v.(string); ok {
			buf.AppendByte('"')
			buf.AppendString(str)
			buf.AppendByte('"')
		} else {
			buf.AppendString(fmt.Sprint(v))
		}
	}

	buf.AppendByte('\n')
	return buf, nil
}
