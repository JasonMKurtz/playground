#!/usr/bin/python
"""
Thu Nov 19 09:55:27 2015 1 10.20.220.39 4609 /share/batch/Exclusion_Group/Archive/20151118-151104-703EXCL_GRP_S200_F4.txt b _ o r ftp_user ftp 0 * c
"""
import operator
import subprocess
log   = "/var/log/xferlog"
range = 10
hosts = dict()
files = dict()
with open(log, 'r') as f:
        for line in f:
                raw = line.strip().split()
                if raw[6] in hosts.keys():
                        hosts[raw[6]] += 1
                else:
                        hosts[raw[6]] = 1

                tlp = '/' + '/'.join(raw[8].split("/")[1:3])
                if tlp in files.keys():
                        files[tlp] += 1
                else:
                        files[tlp] = 1

print "Top {0} FTP connections (by host)".format(range)
for h in sorted(hosts.items(), key=operator.itemgetter(1))[::-1][:range]:
        cmd = "grep {0} {1} | tail -n 1".format(h[0], log)
        p = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
        latest = ' '.join(p.communicate()[0].split()[0:5])
        print "{0}{1}{2}{3}(Latest: {4})".format(h[0], " "*(45 - len(h[0])), h[1], " "*(10 - len(str(h[1]))), latest)
print ""
print "Top {0} FTP connections (by path)".format(range)
for f in sorted(files.items(), key=operator.itemgetter(1))[::-1][:range]:
        cmd = "grep {0} {1} | tail -n 1".format(f[0], log)
        p = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
        latest = ' '.join(p.communicate()[0].split()[0:5])
        print "{0}{1}{2}{3}(Latest: {4})".format(f[0], " "*(45 - len(f[0])), f[1], " "*(10 - len(str(f[1]))), latest)
