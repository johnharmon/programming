#!/bin/bash

rm -f $1
echo "#!/bin/bash" >> $1
echo "" >> $1
echo "sed -n '/^__RPM_START__$/,/^__RPM_END__$/{//!p}' $1 > nano.rpm && truncate -s -1 nano.rpm" >> $1
echo "rpm -i ./nano.rpm" >> $1
echo "code=\$?" >> $1
echo "echo \$code" >> $1
echo "exit \$code" >> $1
echo -e '\n\n__RPM_START__' >> $1 && cat nano-5.6.1-5.el9.x86_64.rpm >> $1 && echo -e  '\n__RPM_END__' >> $1
chmod a+x $1
