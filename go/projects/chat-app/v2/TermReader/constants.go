package main

const (
	STATE_INITIAL_INPUT = iota
	STATE_PARSING_CMD_COUNT
	STATE_CMD_IDENTIFIED
	STATE_PARSING_SUFFIX
	STATE_PENDING_SUFFIX
	STATE_PARSING_MOTION_COUNT
	STATE_PARSING_MOTION
	STATE_PARSING_SPECIAL_SUFFIX
	STATE_EXECUTING_CMD
)

const (
	MR              = "\033[1C"
	ML              = "\033[1D"
	MU              = "\033[1A"
	MD              = "\033[1B"
	MODE_NORMAL     = 0
	MODE_INSERT     = 1
	MODE_VISUAL     = 2
	MODE_CMD        = 3
	CHAR_SPACE      = ' '
	CHAR_EXCLAM     = '!'
	CHAR_QUOTE      = '"'
	CHAR_HASH       = '#'
	CHAR_DOLLAR     = '$'
	CHAR_PERCENT    = '%'
	CHAR_AMPERSAND  = '&'
	CHAR_APOSTRO    = '\''
	CHAR_LPAREN     = '('
	CHAR_RPAREN     = ')'
	CHAR_ASTERISK   = '*'
	CHAR_PLUS       = '+'
	CHAR_COMMA      = ','
	CHAR_MINUS      = '-'
	CHAR_DOT        = '.'
	CHAR_SLASH      = '/'
	CHAR_0          = '0'
	CHAR_1          = '1'
	CHAR_2          = '2'
	CHAR_3          = '3'
	CHAR_4          = '4'
	CHAR_5          = '5'
	CHAR_6          = '6'
	CHAR_7          = '7'
	CHAR_8          = '8'
	CHAR_9          = '9'
	CHAR_COLON      = ':'
	CHAR_SEMICOLON  = ';'
	CHAR_LT         = '<'
	CHAR_EQ         = '='
	CHAR_GT         = '>'
	CHAR_QUESTION   = '?'
	CHAR_AT         = '@'
	CHAR_A          = 'A'
	CHAR_B          = 'B'
	CHAR_C          = 'C'
	CHAR_D          = 'D'
	CHAR_E          = 'E'
	CHAR_F          = 'F'
	CHAR_G          = 'G'
	CHAR_H          = 'H'
	CHAR_I          = 'I'
	CHAR_J          = 'J'
	CHAR_K          = 'K'
	CHAR_L          = 'L'
	CHAR_M          = 'M'
	CHAR_N          = 'N'
	CHAR_O          = 'O'
	CHAR_P          = 'P'
	CHAR_Q          = 'Q'
	CHAR_R          = 'R'
	CHAR_S          = 'S'
	CHAR_T          = 'T'
	CHAR_U          = 'U'
	CHAR_V          = 'V'
	CHAR_W          = 'W'
	CHAR_X          = 'X'
	CHAR_Y          = 'Y'
	CHAR_Z          = 'Z'
	CHAR_LBRACK     = '['
	CHAR_BSLASH     = '\\'
	CHAR_RBRACK     = ']'
	CHAR_CARET      = '^'
	CHAR_UNDERSCORE = '_'
	CHAR_BACKTICK   = '`'
	CHAR_a          = 'a'
	CHAR_b          = 'b'
	CHAR_c          = 'c'
	CHAR_d          = 'd'
	CHAR_e          = 'e'
	CHAR_f          = 'f'
	CHAR_g          = 'g'
	CHAR_h          = 'h'
	CHAR_i          = 'i'
	CHAR_j          = 'j'
	CHAR_k          = 'k'
	CHAR_l          = 'l'
	CHAR_m          = 'm'
	CHAR_n          = 'n'
	CHAR_o          = 'o'
	CHAR_p          = 'p'
	CHAR_q          = 'q'
	CHAR_r          = 'r'
	CHAR_s          = 's'
	CHAR_t          = 't'
	CHAR_u          = 'u'
	CHAR_v          = 'v'
	CHAR_w          = 'w'
	CHAR_x          = 'x'
	CHAR_y          = 'y'
	CHAR_z          = 'z'
	CHAR_LBRACE     = '{'
	CHAR_PIPE       = '|'
	CHAR_RBRACE     = '}'
	CHAR_TILDE      = '~'
)

const (
	CMD_ACCEPT_SUFFIX = 1 << iota
	CMD_ACCEPT_MOTION
)
