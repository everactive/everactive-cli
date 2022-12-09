# VHS documentation
#
# Output:
#   Output <path>.gif               Create a GIF output at the given <path>
#   Output <path>.mp4               Create an MP4 output at the given <path>
#   Output <path>.webm              Create a WebM output at the given <path>
#
# Settings:
#   Set FontSize <number>           Set the font size of the terminal
#   Set FontFamily <string>         Set the font family of the terminal
#   Set Height <number>             Set the height of the terminal
#   Set Width <number>              Set the width of the terminal
#   Set LetterSpacing <float>       Set the font letter spacing (tracking)
#   Set LineHeight <float>          Set the font line height
#   Set Theme <string>              Set the theme of the terminal (JSON)
#   Set Padding <number>            Set the padding of the terminal
#   Set Framerate <number>          Set the framerate of the recording
#   Set PlaybackSpeed <float>       Set the playback speed of the recording
#
# Sleep:
#   Sleep <time>                    Sleep for a set amount of <time> in seconds
#
# Type:
#   Type[@<time>] "<characters>"    Type <characters> into the terminal with a
#                                   <time> delay between each character
#
# Keys:
#   Backspace[@<time>] [number]     Press the Backspace key
#   Down[@<time>] [number]          Press the Down key
#   Enter[@<time>] [number]         Press the Enter key
#   Space[@<time>] [number]         Press the Space key
#   Tab[@<time>] [number]           Press the Tab key
#   Left[@<time>] [number]          Press the Left Arrow key
#   Right[@<time>] [number]         Press the Right Arrow key
#   Up[@<time>] [number]            Press the Up Arrow key
#   Down[@<time>] [number]          Press the Down Arrow key
#   Ctrl+<key>                      Press the Control key + <key> (e.g. Ctrl+C)
#
# Display:
#   Hide                            Hide the subsequent commands from the output
#   Show                            Show the subsequent commands in the output

Output everactive-cli-demo.gif

Set FontSize 16
Set Width 1200
Set Height 600

Type "# Welcome to the Everactive CLI"
Enter 2

# Pause for dramatic effect...
Sleep 500ms


Type "# The Everactive CLI allows you to interact with your Eversensors."
Enter
Sleep 500ms
Type "everactive-cli"
Enter 2
Sleep 2s

Type "# First, initialize the credentials configuration with your API Key:"
Enter
Sleep 500ms
Type "everactive-cli credentials init"
Enter
Hide
Type@25ms "CLIENT_ID"
Enter
Type@25ms "CLIENT_SECRET"
Enter
Show
Sleep 2s

Type "# Now let's get the list of Eversensors."
Enter
Sleep 500ms
Type "everactive-cli list-sensors"
Enter
Sleep 5s

Type "# Get the last reading from one of the sensors."
Enter
Sleep 500ms
Type "everactive-cli data --sensor bc:5e:a1:00:00:00:5e:13"
Enter
Sleep 5s

Type "# Now let's get the last 1 hour of readings."
Enter
Sleep 500ms
Type "everactive-cli data --sensor bc:5e:a1:00:00:00:5e:13 --range 1h"
Enter
Sleep 3s

Type "# Now let's get the last 10 minutes of readings."
Enter
Sleep 500ms
Type "everactive-cli data --sensor bc:5e:a1:00:00:00:5e:13 --range 10m"
Enter
Sleep 3s


Type "# The CLI also supports a streaming mode to get the data as it becomes available."
Enter
Sleep 500ms
Type "everactive-cli stream --sensor bc:5e:a1:00:00:00:5e:13"
Enter
Sleep 10


Type "# Thanks!"
Enter

# Admire the output for a bit.
Sleep 5s


