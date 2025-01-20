# Where should we write the GIF?
Output console.gif

Set Theme "UltraDark"
# Set Margin 10
# Set WindowBar Colorful
# Set Framerate 120

# Typing Speed
Set TypingSpeed 120ms

# Playback speed
# Set PlaybackSpeed 0.9 # a bit slower

# Set up a 1200x600 terminal with 16px font.
Set FontSize 16
Set Width 1200
# Set Width 1000
Set Height 800
# Set Height 660

# run the cexfind console app
Hide
Sleep 100ms
Type "go run ."
Enter
Sleep 2s
Show

# type
Type "lenovo thinkpad x1"
Sleep 1s
Enter
Sleep 2s

# select fourth item
Type @300ms "jjjjj"
Sleep 1s
Type @300ms "jjjjjjj"
Sleep 1s
Enter
Sleep 4s

# type to add lenovo t14s
Tab
Sleep 750ms
Type "; lenovo t14s"
Sleep 500ms
Tab
Sleep 500ms
Type "SW1A 0AA"
Tab
Sleep 500ms
Type "x"
Sleep 1s
Enter
Sleep 2s

# select second item
Type@300ms "jj"
Sleep 2s
Enter
Sleep 7s
