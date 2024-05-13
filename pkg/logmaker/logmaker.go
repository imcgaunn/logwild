package logmaker

import "io"

type LogMakerConfig struct {
	PerSecondRate       int64
	PerMessageSizeBytes int64
	MessageFmt          string
}

func NewWithDefaults() LogMakerConfig {
	return LogMakerConfig{
		PerSecondRate:       1000,
		PerMessageSizeBytes: 1024 * 1000,
		MessageFmt:          "json",
	}
}

func (cfg *LogMakerConfig) StartWriting(*io.Writer w) error {
  return nil
}

func (*LogMakerConfig cfg) StopWritingLogs() error {
  return nil
}
