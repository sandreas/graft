package compare

import (
	"os"

	"github.com/spf13/afero"
	"errors"
	"bytes"
	"math"
)

type Stitch struct {
	BufferSize          int64
	sourceFs            afero.Fs
	source              string
	destinationFs       afero.Fs
	destination         string
	SourceFileStat      os.FileInfo
	DestinationFileStat os.FileInfo

	fi afero.File
	fo afero.File

	isTransferCompleted bool
}

func NewStich(srcFs afero.Fs, src string, dstFs afero.Fs, dst string, bufferSize int64) (*Stitch, error) {
	stitchStrategy := &Stitch{
		BufferSize:    bufferSize,
		sourceFs:      srcFs,
		source:        src,
		destinationFs: dstFs,
		destination:   dst,
	}

	var err error
	stitchStrategy.SourceFileStat, err = srcFs.Stat(src)
	if err != nil {
		return nil, err
	}

	if err := stitchStrategy.initialize(); err != nil {
		return nil, err
	}

	return stitchStrategy, nil
}
func (s *Stitch) initialize() error {
	var err error
	s.DestinationFileStat, err = s.destinationFs.Stat(s.destination)
	if os.IsNotExist(err) {
		s.isTransferCompleted = false
		return nil
	}

	if err != nil {
		return err
	}


	inSize := s.SourceFileStat.Size()
	outSize := s.DestinationFileStat.Size()

	if inSize < outSize {
		return errors.New("source is smaller than destination")
	}

	s.fi, err = s.sourceFs.Open(s.source)
	if err != nil {
		return err
	}
	defer s.fi.Close()

	s.fo, err = s.destinationFs.Open(s.destination)
	if err != nil {
		return err
	}
	defer s.fo.Close()


	if err := s.ensureFileContentsEqual(s.fi, s.fo, 0); err != nil {
		return err
	}

	if outSize <= s.BufferSize {
		s.isTransferCompleted = inSize == outSize
		return nil
	}

	backOffset := outSize-s.BufferSize
	shouldCheckMiddle := true
	if backOffset <= outSize / 2  {
		backOffset = s.BufferSize
		shouldCheckMiddle = false
	}

	if err := s.ensureFileContentsEqual(s.fi, s.fo, backOffset); err != nil {
		return err
	}
	if !shouldCheckMiddle  {
		s.isTransferCompleted = inSize == outSize
		return nil
	}

	middleOffset := int64(math.Floor(float64(outSize / 2) - float64(s.BufferSize / 2)))

	if err := s.ensureFileContentsEqual(s.fi, s.fo, middleOffset); err != nil {
		return err
	}

	s.isTransferCompleted = inSize == outSize
	return nil
}

func (s *Stitch) ensureFileContentsEqual(fi afero.File, fo afero.File, offset int64) error{
	fiBuf := make([]byte, s.BufferSize)
	_, err := fi.ReadAt(fiBuf, offset)

	if err != nil {
		return err
	}

	foBuf := make([]byte, s.BufferSize)
	_, err = fo.ReadAt(foBuf, offset)

	if err != nil {
		return err
	}
	if ! bytes.Equal(fiBuf, foBuf) {
		return errors.New("source file does not match destination file")
	}
	return nil
}



func (s *Stitch) IsComplete() bool {
	return s.isTransferCompleted
}
