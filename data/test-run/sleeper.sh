#!/bin/sh
printf -- "sleeping for $1 seconds... ";
sleep $1;
printf -- "and another $1 seconds... ";
sleep $1;
printf -- "and yet another $1 seconds... ";
sleep $1;
printf -- "done here!\n";
exit 1;
