package resp

import (
	"fmt"
	"io"
)

type Writer struct {
	w io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

func (w *Writer) WriteSimpleString(value string) error {

	_, err := fmt.Fprintf(w.w,"+%s\r\n",value)

	return err
}

func (w *Writer) WriteError(errMsg string) error {

	_, err := fmt.Fprintf(w.w,"-%s\r\n",errMsg)

	return err
}

func (w *Writer) WriteInteger(value int64) error {

	_, err := fmt.Fprintf(w.w,":%d\r\n",value)

	return err
}

func (w *Writer) WriteBulkString(value string) error {

	_, err := fmt.Fprintf(w.w,"$%d\r\n%s\r\n",len(value),value)

	return err
}

func (w *Writer) WriteNullBulkString() error {

	_, err := fmt.Fprint(w.w,"$-1\r\n")

	return err
}

func (w *Writer) WriteArray(values []string,) error {

	fmt.Fprintf(w.w,"*%d\r\n",len(values))

	for _, v := range values {

		fmt.Fprintf(w.w,"$%d\r\n%s\r\n",len(v),v)
	}

	return nil
}