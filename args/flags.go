package args

type flagReg struct {
	flags []flag
}

type flag struct {
	flag    string
	short   string
	fn      func(TokLit, int)
	command bool
}

var FlagReg = &flagReg{flags: make([]flag, 0)}

func (f *flagReg) Register(Flag string, short string, fn func(TokLit, int)) {
	f.flags = append(f.flags, flag{
		flag:  Flag,
		short: short,
		fn:    fn,
	})
}

func (f *flagReg) RegisterCmd(command string, fn func(TokLit, int)) {
	f.flags = append(f.flags, flag{
		flag:    command,
		fn:      fn,
		command: true,
	})
}

// Parse parses a set of tokens.
//
// Returns true if we dispatched to anything
func (f *flagReg) Parse(toks Tokens) bool {
	var didthings bool
	for k, v := range toks {
		switch v.Token {
		case FLAG:
			for _, fo := range f.flags {
				if fo.flag == v.Literal || fo.short == v.Literal {
					fo.fn(v, k)
					didthings = true
					break
				}
			}
		}
	}
	for k, v := range toks {
		if k == 0 && v.Token == STRING { // COMMAND
			for _, fo := range f.flags {
				if v.Literal == fo.flag && fo.command {
					fo.fn(v, k)
					didthings = true
				}
			}
			continue
		}
	}
	return didthings
}
