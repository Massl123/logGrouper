# logGrouper

Parse timebased logs and group them by a common field.  
Built for apache access logs but configurable for others.

## Example
Apache Access Logs
~~~
./logGrouper -i 12h demo-access_log
2004-03-11 Thu
       01:00:00 - 13:00:00                                                                                            63
          GET /cgi-bin/mailgraph.cgi/mailgraph_0_err.png                                                               2
          GET /favicon.ico                                                                                             2
          GET /cgi-bin/mailgraph.cgi/mailgraph_2.png                                                                   2
          GET /twiki/bin/view/Main/WebHome                                                                             2
          GET /cgi-bin/mailgraph2.cgi                                                                                  2
          GET /cgi-bin/mailgraph.cgi/mailgraph_1_err.png                                                               2
          GET /cgi-bin/mailgraph.cgi/mailgraph_1.png                                                                   2
          GET /twiki/pub/TWiki/TWikiLogos/twikiRobot46x50.gif                                                          2
          GET /cgi-bin/mailgraph.cgi/mailgraph_3_err.png                                                               2
          GET /robots.txt                                                                                              2
          GET /cgi-bin/mailgraph.cgi/mailgraph_2_err.png                                                               2
          GET /cgi-bin/mailgraph.cgi/mailgraph_0.png                                                                   2
          GET /cgi-bin/mailgraph.cgi/mailgraph_3.png                                                                   2
          GET /dccstats/stats-hashes.1month.png                                                                        1
          GET /razor.html                                                                                              1
          GET /twiki/bin/view/Main/SpamAssassinTaggingOnly                                                             1
          GET /                                                                                                        1
          GET /twiki/bin/oops/TWiki/Wik                                                                                1
          GET /twiki/bin/view/Main/DCCAndPostFix                                                                       1
          GET /images/image004.jpg                                                                                     1
       13:00:00 - 01:00:00                                                                                           160
          GET /twiki/pub/TWiki/TWikiLogos/twikiRobot46x50.gif                                                         12
          GET /robots.txt                                                                                              4
          GET /twiki/bin/view/Main/WebHome                                                                             4
          GET /razor.html                                                                                              4
          GET /favicon.ico                                                                                             4
          GET /                                                                                                        4
          GET /images/image004.jpg                                                                                     3
          GET /icons/PythonPowered.png                                                                                 3
          GET /icons/gnu-head-tiny.jpg                                                                                 3
          GET /icons/mailman.jpg                                                                                       3
          GET /images/msgops.JPG                                                                                       3
          GET /images/image005.jpg                                                                                     3
          GET /cgi-bin/mailgraph2.cgi                                                                                  2
          POST /mailman/admindb/ppwc                                                                                   2
          GET /ie.htm                                                                                                  2
          GET /twiki/bin/view/Main/WebTopicList                                                                        2
          GET /twiki/bin/view/Main/SpamAssassinAndPostFix                                                              2
          GET /cgi-bin/mailgraph.cgi/mailgraph_1_err.png                                                               2
          GET /cgi-bin/mailgraph.cgi/mailgraph_3.png                                                                   2
          GET /AmavisNew.html                                                                                          2
2004-03-12 Fri
       01:00:00 - 13:00:00                                                                                            29
          GET /robots.txt                                                                                              3
          GET /                                                                                                        2
          GET /SpamAssassin.html                                                                                       1
          GET /dccstats/stats-spam-ratio.1week.png                                                                     1
          GET /dccstats/stats-hashes.1week.png                                                                         1
          GET /dccstats/stats-spam.1month.png                                                                          1
          GET /dccstats/stats-spam-ratio.1month.png                                                                    1
          GET /images/image005.jpg                                                                                     1
          GET /images/image004.jpg                                                                                     1
          HEAD /twiki/bin/view/Main/SpamAssassinDeleting                                                               1
          GET /dccstats/stats-spam.1day.png                                                                            1
          GET /dccstats/stats-hashes.1month.png                                                                        1
          GET /dccstats/stats-hashes.1year.png                                                                         1
          GET /mailman/listinfo/webber                                                                                 1
          GET /twiki/bin/oops/TWiki/1000                                                                               1
          GET /ie.htm                                                                                                  1
          GET /twiki/bin/view/Main/MikeMannix                                                                          1
          GET /dccstats/stats-spam.1week.png                                                                           1
          GET /images/msgops.JPG                                                                                       1
          GET /dccstats/stats-spam-ratio.1day.png                                                                      1
       13:00:00 - 01:00:00                                                                                            69
          GET /twiki/pub/TWiki/TWikiLogos/twikiRobot46x50.gif                                                          5
          GET /                                                                                                        3
          GET /dccstats/stats-hashes.1year.png                                                                         2
          GET /cgi-bin/mailgraph2.cgi                                                                                  2
          GET /cgi-bin/mailgraph.cgi/mailgraph_3_err.png                                                               2
          GET /dccstats/stats-hashes.1month.png                                                                        2
          GET /cgi-bin/mailgraph.cgi/mailgraph_1.png                                                                   2
          GET /dccstats/stats-spam.1month.png                                                                          2
          GET /dccstats/stats-spam-ratio.1year.png                                                                     2
          GET /razor.html                                                                                              2
          GET /cgi-bin/mailgraph.cgi/mailgraph_0.png                                                                   2
          GET /cgi-bin/mailgraph.cgi/mailgraph_2_err.png                                                               2
          GET /dccstats/stats-spam.1week.png                                                                           2
          GET /dccstats/stats-spam-ratio.1month.png                                                                    2
          GET /twiki/bin/view/Main/WebHome                                                                             2
          GET /twiki/bin/view/Main/SpamAssassinUsingRazorAndDCC                                                        2
          GET /dccstats/stats-hashes.1week.png                                                                         2
          GET /dccstats/stats-spam.1year.png                                                                           2
          GET /dccstats/stats-spam-ratio.1day.png                                                                      2
          GET /cgi-bin/mailgraph.cgi/mailgraph_0_err.png                                                               2


Interval: 12h0m0s, Output timezone: +0100 CET, Unmatched Lines: 1
~~~

## Usage
~~~
Usage: logGrouper [options] file1 [file2 file...]
Use file name "-" to read from stdin.

Group lines in log by time and occurance.

Copyright (c) 2020 Marcel Freundl <github.com/Massl123>

Profiles:
Name                                     Log format                                                                                           Time format              
apacheAccessLog-first-group              ^.*?\[(?P<timestamp>.+?)\].*?"(?P<group>.+ [/\*].*?[/?\ ]).+?" .*$                                   2/Jan/2006:15:04:05 -0700
apacheAccessLog-full                     ^.*?\[(?P<timestamp>.+?)\].*?"(?P<group>.+? .+?) .+?" .*$                                            2/Jan/2006:15:04:05 -0700

Options:
  -f, --format string       LogFormat regexp in GoLang Format. Match group "timestamp" and "group" have to exist. See https://golang.org/pkg/regexp/syntax/.
  -i, --interval string     Interval to group by. Format like 15m (Units supported: ns, us (or µs), ms, s, m, h). (default "15m")
  -l, --limit int           Limit output per timeslot. (default 20)
  -p, --profile string      Profile to use. Profile loads LogFormat and TimeFormat. Use LogFormat and TimeFormat parameters to override. (default "apacheAccessLog-full")
  -t, --timeFormat string   Time format for "timestamp" match group. Given in GoLang format, see https://golang.org/pkg/time/#Parse.
  -v, --verbose             Verbose output (show unparsed lines)
~~~

## TODOs
* Add --start-time and --end-time to filter logdata
