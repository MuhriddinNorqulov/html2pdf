package converter

import (
	"context"
	"io"
	"os/exec"
	"strings"
	"time"
)

type PdfConverter struct{}

// @inject
func NewPdfConverter() *PdfConverter { return &PdfConverter{} }

// ConvertStream: HTML (string) -> wkhtmltopdf stdin, PDF -> stdout (stream)
// Caller MUST Close() returned reader (masalan: defer r.Close()).
func (c *PdfConverter) ConvertStream(
	ctx context.Context,
	content string,
	params map[string]string,
) (io.ReadCloser, error) {

	// timeout bo‘lmasa ham ishlaydi, lekin yaxshi amaliyot:
	if ctx == nil {
		ctx = context.Background()
	}
	// xohlasangiz tashqaridan timeout bering; bo‘lmasa default:
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)

	args := make([]string, 0, len(params)*2+2)
	for k, v := range params {
		key := strings.TrimLeft(k, "-") // "--margin-top" -> "margin-top"
		args = append(args, "--"+key)
		if v != "" {
			args = append(args, v)
		}
	}
	args = append(args, "-", "-") // stdin -> stdout

	cmd := exec.CommandContext(ctx, "wkhtmltopdf", args...)
	cmd.Stdin = strings.NewReader(content)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return nil, err
	}
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, err
	}

	// stderr pipe to'lib qolib deadlock bo'lmasin
	go io.Copy(io.Discard, stderr)

	// stdout reader + cmd lifecycle ni bitta ReadCloser qilib qaytaramiz
	return &wkStream{
		r:      stdout,
		cmd:    cmd,
		cancel: cancel,
	}, nil
}

type wkStream struct {
	r      io.ReadCloser
	cmd    *exec.Cmd
	cancel context.CancelFunc
	once   bool
}

func (s *wkStream) Read(p []byte) (int, error) {
	n, err := s.r.Read(p)
	if err == io.EOF {
		// process tugaganini tekshiramiz (exit code)
		_ = s.cmd.Wait()
		s.Close() // cancel ham bo'ladi
	}
	return n, err
}

func (s *wkStream) Close() error {
	if s.once {
		return nil
	}
	s.once = true

	// context cancel (deadline/timeout)
	if s.cancel != nil {
		s.cancel()
	}

	// client uzib yuborsa processni to'xtatamiz
	if s.cmd != nil && s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
	}

	_ = s.r.Close()
	_ = s.cmd.Wait()
	return nil
}
