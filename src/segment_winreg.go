package main

import (
	"fmt"
)

type winreg struct {
	props properties
	env   environmentInfo

	Value string
}

const (
	// path from the supplied root under which the key exists
	RegistryPath Property = "path"
	// key within full reg path formed from two above
	RegistryKey Property = "key"
	// Fallback is the text to display if the key is not found
	Fallback Property = "fallback"
)

func (wr *winreg) init(props properties, env environmentInfo) {
	wr.props = props
	wr.env = env
}

func (wr *winreg) enabled() bool {
	if wr.env.getRuntimeGOOS() != windowsPlatform {
		return false
	}

	registryPath := wr.props.getString(RegistryPath, "")
	registryKey := wr.props.getString(RegistryKey, "")
	fallback := wr.props.getString(Fallback, "")

	var regValue *windowsRegistryValue
	regValue, _ = wr.env.getWindowsRegistryKeyValue(registryPath, registryKey)

	if regValue != nil {
		switch regValue.valueType {
		case regString:
			wr.Value = regValue.str
			return true
		case regDword:
			wr.Value = fmt.Sprintf("0x%08X", regValue.dword)
			return true
		case regQword:
			wr.Value = fmt.Sprintf("0x%016X", regValue.qword)
			return true
		}
	}

	if len(fallback) > 0 {
		wr.Value = fallback
		return true
	}

	return false
}

func (wr *winreg) string() string {
	segmentTemplate := wr.props.getString(SegmentTemplate, "{{ .Value }}")
	return wr.templateString(segmentTemplate)
}

func (wr *winreg) templateString(segmentTemplate string) string {
	template := &textTemplate{
		Template: segmentTemplate,
		Context:  wr,
		Env:      wr.env,
	}
	text, err := template.render()
	if err != nil {
		return err.Error()
	}
	return text
}
