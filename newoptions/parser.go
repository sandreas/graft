package newoptions

type BitFlag byte


type BitFlagParser struct {
	activeFlags BitFlag
}


func NewBitFlagParser( params ...BitFlag) *BitFlagParser {
	parser := &BitFlagParser{}
	parser.parse(params)
	return parser
}
func (parser *BitFlagParser) parse(params []BitFlag) {
	size := len(params)

	parser.activeFlags = 0x00
	for i := 0; i < size; i++ {
		parser.activeFlags |= params[i]
	}
}

func (parser *BitFlagParser) hasFlag(flagToCheck BitFlag) bool {
	return parser.activeFlags & flagToCheck != 0
}