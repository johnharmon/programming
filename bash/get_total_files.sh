 #!/bin/bash


 ls -l -R / 2>/dev/null | sed '/\(^$\)\|\(^total\)\|\(\/\)/d' | wc -l 

