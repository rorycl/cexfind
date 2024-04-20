# charm vhs screen recording file for the cli app
#
# Where should we write the GIF?
Output cli.gif

# Set Theme "UltraDark"
# Set Margin 10
# Set WindowBar Colorful
# Set Framerate 120

# set the PS1 prompt to '> '!

# Typing Speed
Set TypingSpeed 100ms

# Default is to set up a 1200x600 terminal with 16px font.
Set FontSize 16
Set Width 1000
Set Height 660

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
Sleep 3s

# clear
Hide
Type "clear"
Enter
Show

# select item
Type "./cli -query 'lenovo t480s' | head -n 20"
Enter
Sleep 3s

# clear
Hide
Type "clear"
Enter
Show

# select more items
Type "./cli -query 'lenovo t480s' -query 'lenovo x390 yoga' -strict"
Enter
Sleep 3s

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
