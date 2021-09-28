package utils

import (
	"fmt"
	"os"

	"github.ibm.com/Gufran-Baig/fargo-fb-poc/api/apiproto"
)

type Writer struct {
	fw map[string]*os.File
}

func New() *Writer {
	return &Writer{
		fw: make(map[string]*os.File),
	}
}

func (w *Writer) Write(event *apiproto.Event, decrypt bool) error {
	_, ok := w.fw[event.AccessKey]
	if !ok {
		err := w.initfile(event.AccessKey)
		if err != nil {
			return err
		}
	}

	err := w.writefile(event, decrypt)
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) Close() {
	for _, f := range w.fw {
		_ = f.Close()
	}
}

func (w *Writer) initfile(accessKey string) error {
	f, err := os.OpenFile(fmt.Sprintf("../%s.log", accessKey), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	w.fw[accessKey] = f

	return nil
}

func (w *Writer) writefile(event *apiproto.Event, decrypt bool) error {
	f, _ := w.fw[event.AccessKey]
	e, err := MarshalEvent(event, decrypt)
	if err != nil {
		return fmt.Errorf("Failed to write event with error %v", err)
	}

	if _, err = f.WriteString(e); err != nil {
		return fmt.Errorf("Failed to write event with error %v", err)
	}

	return nil
}
