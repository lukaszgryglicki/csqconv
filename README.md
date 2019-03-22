# csqconv

Convert FLIR camera `*.csq` files into a `*.pp4` H.264 video file (you can optionaly save all intermediate files: `*.raw`, `*.jpegls`, `*.png` etc.)

# Requirements

- `ffmpeg`.
- [jpegbw](https://github.com/lukaszgryglicki/jpegbw).

# Usage

- `make && make install`.
- `RC=1 GC=1 BC=1 RLO=55 RHI=3 GLO=25 GHI=25 BLO=3 BHI=55 NA=1 LIB=libjpegbw.so RF="saturate(x1, .0001_.0001, .9999_1)" GF="saturate(x1, .0001_.0001, .9999_.3)" BF="saturate(x1, .0001_.0001, .9999_.2)" csqconv filename.csq`.
- `RLO=3 RHI=3 GLO=3 GHI=3 BLO=3 BHI=3 NA=1 csqconv ./small.csq`.
- C optimized rainbow from grayscale conversion: `RLO=5 RHI=5 GLO=5 GHI=5 BLO=5 BHI=5 NA=1 LIB=libjpegbw.so RF="gsrainbowr(x1)" GF="gsrainbowg(x1)" BF="gsrainbowb(x1)" csqconv small.csq`.
- `RC=1 GC=1 BC=1 RLO=3 RHI=3 GLO=3 GHI=3 BLO=3 BHI=3 NA=1 LIB=libjpegbw.so RF="gsrainbowre(x1, .02)" GF="gsrainbowge(x1, .02)" BF="gsrainbowbe(x1, .02)" csqconv small.csq`.
- `RC=1 GC=1 BC=1 RLO=3 RHI=3 GLO=3 GHI=3 BLO=3 BHI=3 NA=1 LIB=libjpegbw.so RF="(3*x1+gsrainbowr(x1))/4" GF="(3*x1+gsrainbowg(x1))/4" BF="(3*x1+gsrainbowb(x1))/4" csqconv FLIR0009.csq`.
- Alarms: `RC=1 GC=1 BC=1 RLO=1 RHI=1 GLO=1 GHI=1 BLO=1 BHI=1 NA=1 LIB=libjpegbw.so BF="saturate(x1, .02_1,.98)" GF="saturate(x1, .02, .98)" RF="saturate(x1, 0, .98_1)" csqconv f.csq`
- Alarm and rainbow combined: `RC=1 GC=1 BC=1 RLO=1 RHI=1 GLO=1 GHI=1 BLO=1 BHI=1 NA=1 LIB=libjpegbw.so BF="gsrainbowb(saturate(x1, .01_1,.99))" GF="gsrainbowg(saturate(x1, .01, .99))" RF="gsrainbowr(saturate(x1, 0, .99_1))" csqconv f.csq`.
- Alarm, rainbow and white-hot/black-cold: `RC=1 GC=1 BC=1 RLO=1 RHI=1 GLO=1 GHI=1 BLO=1 BHI=1 NA=1 LIB=libjpegbw.so BF="(3*x1+2*gsrainbowb(saturate(x1, .01_1,.99)))/5" GF="(3*x1+2*gsrainbowg(saturate(x1, .01, .99)))/5" RF="(3*x1+2*gsrainbowr(saturate(x1, 0, .99_1)))/5" csqconv f.csq`.
- See `jpegbw` README.md to see multiple options that can be passed to this tool.
