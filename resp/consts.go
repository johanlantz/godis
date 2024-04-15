package resp

const (
	DT_SIMPLE_STRING    = '+'
	DT_SIMPLE_ERROR     = '-'
	DT_INTEGER          = ':'
	DT_BULK_STRINGS     = '$'
	DT_ARRAYS           = '*'
	DT_NULLS            = '_'
	DT_BOOLEANS         = '*'
	DT_DOUBLES          = ','
	DT_BIG_NUMBERS      = '('
	DT_BULK_ERRORS      = '|'
	DT_VERBATIM_STRINGS = '='
	DT_MAPS             = '%'
	DT_SETS             = '~'
	DT_PUSHES           = '>'
)

const (
	RESP_OK  = "OK"
	RESP_ERR = "ERR"
)

var allCommands = []byte{DT_SIMPLE_STRING, DT_SIMPLE_ERROR, DT_INTEGER, DT_BULK_STRINGS, DT_ARRAYS, DT_NULLS, DT_BOOLEANS, DT_DOUBLES, DT_BIG_NUMBERS, DT_BULK_ERRORS, DT_VERBATIM_STRINGS, DT_MAPS, DT_SETS, DT_PUSHES}
