#!/usr/bin/python
# -*- coding:utf-8 -*-
# Author:      LiuSha
from __future__ import unicode_literals
from __future__ import print_function

from datetime import datetime
from time import time as unix_time
from optparse import OptionParser
from subprocess import PIPE, Popen, STDOUT

import os


def shell_call(cmd, output=True):
    if output:
        process = Popen(cmd, stdout=PIPE, shell=True, universal_newlines=True, stderr=STDOUT)
        output, _ = process.communicate()
        status = process.poll()

        if output[-1:] == '\n':
            output = output[:-1]

        return status, output.decode('utf-8')
    else:
        os.system('/usr/bin/nohup %s &' % cmd)


def remove_log_file(path, filename):
    filename = os.path.join(path, filename)
    if not os.path.exists(filename):
        raise IOError('No such file or directory: %s' % filename)

    os.remove(filename)


def parser_args_check(args):
    _opt, _ = args()

    if not _opt.dir or not _opt.bak or not _opt.pid:
        raise OSError('script options requires')

    if not os.path.exists(_opt.dir) or not os.path.exists(_opt.bak):
        raise IOError('No such file or directory: %s or %s' % (_opt.bak, _opt.dir))

    if not os.path.isfile(_opt.pid):
        raise IOError('No such file or directory: %s' % _opt.pid)

    return _opt


def main():
    usage = "<your script> [options]"
    parser = OptionParser(usage)
    parser.add_option("--dir", help="Set the log source path")
    parser.add_option("--key", help="Set keyword matching log")
    parser.add_option("--bak", help="Set the log backup path")
    parser.add_option("--pid", help="Set the nginx pid path")
    parser.add_option("--days", type=int, help="Backup How many days, 20 days default")
    options = parser_args_check(parser.parse_args)

    today = datetime.now().strftime('_%Y-%m-%d')
    if options.days is None:
        del_date = (unix_time() - (20 * 86400))
    else:
        del_date = (unix_time() - (options.days * 86400))

    b_file_list = os.listdir(options.dir)
    backup_dir = os.path.join(options.bak, 'backup_log%s' % today)
    if options.key:
        b_file_list = [item for item in b_file_list if options.key in item]

    if os.path.exists(backup_dir):
        exit('Backup already exists today!')

    print('Start Working[%s]' % datetime.now())
    if not os.path.exists(backup_dir) and os.path.exists(os.path.dirname(backup_dir)):
        os.mkdir(backup_dir)

    for item in b_file_list:
        shell_call('/bin/mv {src} {des}'.format(
            src=os.path.join(options.dir, item),
            des=os.path.join(backup_dir, item)
        ) + today)
        print('Backup Files: %s' % os.path.join(backup_dir, item))

    d_dir_list = [item for item
                  in os.listdir(options.bak)
                  if 'backup_log' in item]

    for item in d_dir_list:
        if os.stat(os.path.join(options.bak, item)).st_mtime < del_date:
            shell_call('/bin/rm -rf %s' % os.path.join(options.bak, item))
            print('Removing Backup Files: %s' % os.path.join(options.bak, item))

    if os.path.exists(options.pid):
        shell_call('kill -USR1 `cat {pid}`'.format(pid=options.pid))
        print('Reload Nginx!\n')

if __name__ == '__main__':
    main()

