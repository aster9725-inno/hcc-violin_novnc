package encodings

import (
	"fmt"
	"hcc/violin-novnc/common"
	"hcc/violin-novnc/lib/logger"
	"io"
)

type TightPngEncoding struct {
	bytes []byte
}

func (z *TightPngEncoding) WriteTo(w io.Writer) (n int, err error) {
	return w.Write(z.bytes)
}

func (*TightPngEncoding) Type() int32 { return int32(common.EncTightPng) }

func (t *TightPngEncoding) Read(pixelFmt *common.PixelFormat, rect *common.Rectangle, r *common.RfbReadHelper) (common.IEncoding, error) {
	bytesPixel := calcTightBytePerPixel(pixelFmt)
	r.StartByteCollection()
	defer func() {
		t.bytes = r.EndByteCollection()
	}()

	//var subencoding uint8
	compctl, err := r.ReadUint8()
	if err != nil {
		logger.Logger.Printf("error in handling tight encoding: %v", err)
		return nil, err
	}
	logger.Logger.Printf("bytesPixel= %d, subencoding= %d", bytesPixel, compctl)

	//move it to position (remove zlib flush commands)
	compType := compctl >> 4 & 0x0F

	logger.Logger.Printf("afterSHL:%d", compType)
	switch compType {
	case TightPNG:
		len, err := r.ReadCompactLen()
		_, err = r.ReadBytes(len)

		if err != nil {
			return t, err
		}

	case TightFill:
		r.ReadBytes(int(bytesPixel))
	default:
		return nil, fmt.Errorf("unknown tight compression %d", compType)
	}
	return t, nil
}
