# augur — OpenTelemetry Collector config linter
# https://github.com/starkross/augur

Output demo.gif

Set Shell zsh
Set FontSize 16
Set Width 1200
Set Height 800
Set Padding 20
Set Theme "Catppuccin Mocha"
Set TypingSpeed 50ms

Type "cat examples/bad.yaml"
Enter
Sleep 5s

Type "./augur examples/bad.yaml"
Enter
Sleep 8s

Type "./augur examples/good.yaml"
Enter
Sleep 6s

Type "./augur -o json examples/bad.yaml | head -20"
Enter
Sleep 6s

Type "./augur --strict examples/good.yaml"
Enter
Sleep 6s

Type "./augur --skip OTEL-010,OTEL-015 examples/bad.yaml"
Enter
Sleep 8s
