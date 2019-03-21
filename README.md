# csqconv

Convert FLIR camera `*.csq` files into a `*.pp4` H.264 video file (you can optionaly save all intermediate files: *.raw, *.jpegls, *.png etc.)

# Requirements

- `ffmpeg`.
- [jpegbw](https://github.com/lukaszgryglicki/jpegbw).

# Usage

- `make && make install`.
- `RLO=55 RHI=3 GLO=25 GHI=25 BLO=3 BHI=55 NA=1 LIB=libjpegbw.so RF="saturate(x1, .0001_.0001, .9999_1)" GF="saturate(x1, .0001_.0001, .9999_.3)" BF="saturate(x1, .0001_.0001, .9999_.2)" csqconv filename.csq`.
- `RLO=3 RHI=3 GLO=3 GHI=3 BLO=3 BHI=3 NA=1 csqconv ./small.csq`.
- C optimized rainbow from grayscale conversion: `RLO=5 RHI=5 GLO=5 GHI=5 BLO=5 BHI=5 NA=1 LIB=libjpegbw.so RF="gsrainbowr(x1)" GF="gsrainbowg(x1)" BF="gsrainbowb(x1)" csqconv small.csq`.
- `RLO=3 RHI=3 GLO=3 GHI=3 BLO=3 BHI=3 NA=1 LIB=libjpegbw.so RF="gsrainbowre(x1, .02)" GF="gsrainbowge(x1, .02)" BF="gsrainbowbe(x1, .02)" csqconv small.csq`.
- See `jpegbw` REAME.md to see multiple options that can be passed to this tool.
