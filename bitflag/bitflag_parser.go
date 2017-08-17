package bitflag


type Flag byte


type Parser struct {
	activeFlags Flag
}


func NewParser( params ...Flag) *Parser {
	parser := &Parser{}
	parser.parse(params)
	return parser
}
func (parser *Parser) parse(params []Flag) {
	size := len(params)

	parser.activeFlags = 0x00
	for i := 0; i < size; i++ {
		parser.activeFlags |= params[i]
	}
}

func (parser *Parser) SetFlag(flagToSet Flag) {
	parser.activeFlags |= flagToSet
}

func (parser *Parser) HasFlag(flagToCheck Flag) bool {
	return parser.activeFlags & flagToCheck != 0
}
