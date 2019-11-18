# Godstat

Use Go to port all the features of [dstat](https://github.com/dagwieers/dstat).

## How to use

```bash
# git clone https://github.com/JasonJe/godstat.git
# go run dstat.go
```

## 
Get system basic information

```bash
go run dstat.go --info

```

## Command line argument

```bash
# go run dstat.go --help
Usage of godstat:
      --aio                             enable aio stats.
  -c, --cpu                             enable cpu stats.
  -C, --cpuarray strings                example: 0,3,total (default [total])
      --delay int                       time delay. (default 1)
  -d, --disk                            enable disk stats.
  -D, --diskarray strings               example: total,hda (default [total])
  -T, --epoch                           enable time counter. (Seconds since epoch)
  -f, --filesystem                      enable filesystem stats.
  -h, --help                            help
  -i, --info                            show system information.
  -r, --io                              enable io stats.
                                                (I/O requests completed)
      --ipc                             enable ipc stats.
  -l, --load                            enable load stats.
      --lock                            enable lock stats.
  -m, --mem                             enable memory stats.
  -n, --net                             enable net stats.
  -N, --netarray strings                example: eth1,total (default [total])
  -o, --out string[="2019-11-18.csv"]   write CSV output to file. example: --out=./out.csv
  -g, --page                            enable page stats.
  -p, --proc                            enable process stats.
      --raw                             enable raw  stats.
      --socket                          enable socket stats.
  -s, --swap                            enable swap stats.
  -S, --swaparray strings               example: swap1,total (default [total])
  -y, --sys                             enable system stats.
      --tcp                             enable tcp stats.
  -t, --time                            enable time/date output.
      --udp                             enable udp stats.
      --unix                            enable unix stats.
      --vm                              enable vm stats.
      --zones                           enable zoneinfo stats.
```

