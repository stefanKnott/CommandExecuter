# CommandExecuter

Files included:
  blog1-3.txt: Example text used for testing
  commandExecuter.go:  Here lies the Command Executer program.  Performs checksum, word count, and word frequency.
  command_file.txt: Original command file
  command_file_invalid.txt: Command file with invalid commands, used for testing
  command_file_spaced.txt: Command file with commands surrounded by varying white space, used for testing
  command_file_mixed.txt: Command file with commands surrounded by varying white space, and varying capitalizations of command names, used for testing
  command_test.go: File used for displaying output of commandExecuter running various files
    Note: this does not perform unit tests, but is just used for gauging valid program performace via output to stdout
  

To build and run:
  go build commandExecuter.go
  ./commandExecuter <command_file_name>

To display test script output:
  go test -v

Done as programming assignment for interview
