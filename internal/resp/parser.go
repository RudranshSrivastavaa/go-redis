package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Parser struct {
	reader *bufio.Reader
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(r),
	}
}

func (p *Parser) readLine() (string, error) {

	line, err := p.reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")

	return line, nil
}

func (p *Parser) Parse() ([]string, error) {

	firstByte, err := p.reader.Peek(1)

	if err != nil {
		return nil, err
	}

	if firstByte[0] == ArrayPrefix {
		return p.ParseCommand()
	}

	line, err := p.readLine()

	if err != nil {
		return nil, err
	}

	return ParseInline(line), nil
}

func (p *Parser) ParseCommand() ([]string, error) {

	prefix, err := p.reader.ReadByte()

	if err != nil {
		return nil, err
	}

	if prefix != ArrayPrefix {

		return nil, fmt.Errorf(
			"expected array prefix '*'",
		)
	}

	countLine, err := p.readLine()

	if err != nil {
		return nil, err
	}

	count, err := strconv.Atoi(countLine)

	if err != nil {
		return nil, err
	}

	args := make([]string, 0, count)

	for i := 0; i < count; i++ {

		arg, err := p.parseBulkString()

		if err != nil {
			return nil, err
		}

		args = append(args, arg)
	}

	return args, nil
}
func (p *Parser) parseBulkString() (string, error) {

	prefix, err := p.reader.ReadByte()

	if err != nil {
		return "", err
	}

	if prefix != BulkStringPrefix {

		return "", fmt.Errorf(
			"expected bulk string '$'",
		)
	}

	lengthLine, err := p.readLine()

	if err != nil {
		return "", err
	}

	length, err := strconv.Atoi(lengthLine)

	if err != nil {
		return "", err
	}

	buf := make([]byte, length)

	_, err = io.ReadFull(
		p.reader,
		buf,
	)

	if err != nil {
		return "", err
	}

	crlf := make([]byte, 2)

	_, err = io.ReadFull(
		p.reader,
		crlf,
	)

	if err != nil {
		return "", err
	}

	return string(buf), nil
}