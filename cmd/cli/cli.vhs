# charm vhs screen recording file for the cli app
#
# command reference: https://github.com/charmbracelet/vhs?tab=readme-ov-file#vhs-command-reference
#
# Where should we write the GIF?
Output cli.gif

# Set Theme "UltraDark"
# Set Margin 10
# Set WindowBar Colorful
# Set Framerate 120

# Require the cli binary!
Require "./cli"

# Typing Speed
Set TypingSpeed 120ms

# Default is to set up a 1200x600 terminal with 16px font.
Set FontSize 16
Set Width 1200
Set Height 792

# set the PS1 prompt to '> '!
# run the cli in a subshell with a custom prompt
# https://wiki.archlinux.org/title/Bash/Prompt_customization
Hide
Type "bash"
Enter

Type `RESET="\[$(tput sgr0)\]"`
Enter
Type `CYAN="\[$(tput setaf 6)\]"`
Enter
Type `export PS1="${CYAN}>${RESET} "`
Enter
Type "clear"
Enter
Sleep 100ms
Show

# show help
Type "./cli -h"
Enter
Sleep 4s

# clear
Hide
Type "clear"
Enter
Show

# select item
Hide
Type "# a basic query"
# Ctrl+D
Enter
Show
Sleep 500ms

Type "./cli -query 'Lenovo t14 Gen4' | head -n 19"
Enter
Sleep 4s

# clear
Hide
Type "clear"
Enter
Show

# select more items
Hide
Type "# only show strict matches"
Enter
Show
Sleep 500ms

Type "./cli -query 'Lenovo t14 Gen4' -strict"
Enter
Sleep 4s

# clear
Hide
Type "clear"
Enter
Show
Sleep 500ms

# show buying prices, stores
Hide
Type "# also show cash/exchange pricing and stores"
Enter
Show
Sleep 500ms

Type "./cli -query 'Lenovo t14 Gen4' -strict -verbose"
Enter
Sleep 4s

# clear
Hide
Type "clear"
Enter
Show
Sleep 500ms

# show buying prices, stores and distance to nearest stores
Hide
Type "# also show distance to nearest stores (postcode implies verbose)"
Enter
Show
Sleep 500ms

Type `./cli -query 'Lenovo t14 Gen4' -strict -postcode "B3 2BJ"`
Enter
Sleep 4s

# clear
Hide
Type "clear"
Enter
Show
Sleep 500ms

# and do multiple queries
Hide
Type "# ... and you can also do multiple queries at once"
Enter
Show
Sleep 500ms

Type `./cli -query 'Lenovo t14 Gen4' -query "lenovo x390" -strict -postcode "B3 2BJ"`
Enter
Sleep 6s

# -query "lenovo x390"
# # type lenovo t480s
# Tab
# Sleep 750ms
# Type ", lenovo x390"
# Sleep 1s
# Tab
# Sleep 1s
# Type "x"
# Sleep 1s
# Enter
# Sleep 2s
# 
# # select second item
# Type "j"
# Sleep 1s
# Type "j"
# Sleep 3s
