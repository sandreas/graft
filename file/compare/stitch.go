package compare

import (
	"os"

	"github.com/spf13/afero"
	"errors"
	"bytes"
	"math"
)

type Stitch struct {
	BufferSize             int64
	SourceFileStat         os.FileInfo
	DestinationFileStat    os.FileInfo
	SourceFilePointer      afero.File
	DestinationFilePointer afero.File

	isTransferCompleted bool
}

func NewStich(src afero.File, dst afero.File, bufferSize int64) (*Stitch, error) {
	stitchStrategy := &Stitch{
		BufferSize:             bufferSize,
		SourceFilePointer:      src,
		DestinationFilePointer: dst,
	}

	var err error
	stitchStrategy.SourceFileStat, err = src.Stat()
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
	s.DestinationFileStat, err = s.DestinationFilePointer.Stat()
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

	if outSize == 0 {
		return nil
	}

	if s.BufferSize > outSize  {
		s.BufferSize = outSize
	}

	if err := s.ensureFileContentsEqual(0); err != nil {
		return err
	}

	if outSize <= s.BufferSize {
		s.isTransferCompleted = inSize == outSize
		return nil
	}

	backOffset := outSize - s.BufferSize
	shouldCheckMiddle := true
	if backOffset <= outSize/2 {
		backOffset = s.BufferSize
		shouldCheckMiddle = false
	}

	if err := s.ensureFileContentsEqual(backOffset); err != nil {
		return err
	}
	if !shouldCheckMiddle {
		s.isTransferCompleted = inSize == outSize
		return nil
	}

	middleOffset := int64(math.Floor(float64(outSize/2) - float64(s.BufferSize/2)))

	if err := s.ensureFileContentsEqual(middleOffset); err != nil {
		return err
	}

	s.isTransferCompleted = inSize == outSize
	return nil
}

func (s *Stitch) ensureFileContentsEqual(offset int64) error {
	fiBuf := make([]byte, s.BufferSize)
	_, err := s.SourceFilePointer.ReadAt(fiBuf, offset)

	if err != nil {
		return err
	}

	foBuf := make([]byte, s.BufferSize)
	_, err = s.DestinationFilePointer.ReadAt(foBuf, offset)

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
