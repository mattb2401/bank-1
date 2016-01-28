#! /bin/bash

openssl s_client -showcerts -debug -connect thebankoftoday.com:3300 -no_ssl2 -bugs
