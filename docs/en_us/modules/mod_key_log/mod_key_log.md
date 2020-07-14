# mod_key_log

## Introduction 

mod_key_log writes tls key logs in NSS key log format so that external
programs(eg. wireshark) can decrypt TLS connections for trouble shooting.

For more information about NSS key log format, see:
https://developer.mozilla.org/en-US/docs/Mozilla/Projects/NSS/Key_Log_Format

## Module Configuration

### Description
conf/mod_key_log/mod_key_log.conf

| Config Item | Description                             |
| ----------- | --------------------------------------- |
| Log.LogPrefix | String<br>Filename prefix for log |
| Log.LogDir | String<br>Directory of log files |
| Log.RotateWhen | String<br>Inteval to rotate log file |
| Log.BackupCount | Integer<br>Max number of rotated log files |

### Example
```ini
[Log]
# filename prefix for log 
LogPrefix = key

# log directory 
LogDir = ../log

# interval to rotate logs: M/H/D
# - M: minute
# - H: hour
# - D: day
RotateWhen = H 

# max number of rotated log files
BackupCount = 3
```

