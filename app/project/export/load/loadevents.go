package load

import (
	"encoding/json/jsontext"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/util"
)

func LoadEvents(exportPath string, fs filesystem.FileLoader, logger util.Logger) (map[string]jsontext.Value, model.Events, error) {
	eventsPath := util.StringFilePath(exportPath, "events")
	if !fs.IsDir(eventsPath) {
		return nil, nil, nil
	}
	eventNames := fs.ListJSON(eventsPath, nil, false, logger)
	events := make(model.Events, 0, len(eventNames))
	eventFiles := make(map[string]jsontext.Value, len(eventNames))
	for _, eventName := range eventNames {
		fn := util.StringFilePath(eventsPath, eventName)
		content, err := fs.ReadFile(fn)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to read export event file from [%s]", fn)
		}
		m, err := util.FromJSONObj[*model.Event](content)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to read export event JSON from [%s]", fn)
		}
		eventFiles[m.Name] = content
		events = append(events, m)
	}
	return eventFiles, events, nil
}
