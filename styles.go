package click

type RenderStyle struct {
	Indent            string
	IndentLevel       int
	ClauseNamePrefix  string
	ClauseNameSuffix  string
	ArgumentPrefix    string
	ArgumentSuffix    string
	ArgumentDelimiter string
}

var (
	defaultStyle = RenderStyle{
		ClauseNamePrefix:  " ",
		ClauseNameSuffix:  " ",
		ArgumentPrefix:    "",
		ArgumentSuffix:    "",
		ArgumentDelimiter: ", ",
	}
	prettyStyle = RenderStyle{
		Indent:            "\t",
		ClauseNamePrefix:  "",
		ClauseNameSuffix:  "\n",
		ArgumentPrefix:    "\t",
		ArgumentSuffix:    "\n",
		ArgumentDelimiter: ",",
	}
)
