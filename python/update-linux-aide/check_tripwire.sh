#!/bin/bash
tripwire_path='/ur/local/sbin/tripwire'
last_tripwire_report_path = '/usr/local/lib/tripwire/report/'
last_tripwire_report=$(ls -ltr ${last_tripwire_report_path}/*.twr | tail -n 1 | awk '{print $8}')
printf "Last tripwire report found at ${last_tripwire_report_path}${last_tripwire_report}\n"
$tripwire_path -m u -r ${tripwire_report_path}${last_tripwire_report}
$tripwire_path --check >> ${last_tripwire_report_path}PostDatabaseUpdate_$(hostname)_$(date +"%Y-%m-%d").log
printf "Tripwire check log written to ${last_tripwire_report_path}PostDatabaseUpdate_$(hostname)_$(date +"%Y-%m-%d").log\n"
printf "Displaying last 20 lines of the tripwire log, please check this output to ensure that a line saying \"No Errors\" is present: \n\n"
tail -n 20  ${last_tripwire_report_path}PostDatabaseUpdate_$(hostname)_$(date +"%Y-%m-%d").log


