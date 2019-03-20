# csqconv

Convert FLIR camera `*.csq` files into a sequence of loseless JPEG's and a video file readable by standard tools.

# Requirements

- `ffmpeg`.
- [jpegbw](https://github.com/lukaszgryglicki/jpegbw).

# Usage

- `make && make install`.
- `RLO=55 RHI=3 GLO=25 GHI=25 BLO=3 BHI=55 NA=1 LIB=libjpegbw.so RF="saturate(x1, .0001_.0001, .9999_1)" GF="saturate(x1, .0001_.0001, .9999_.3)" BF="saturate(x1, .0001_.0001, .9999_.2)" csqconv filename.csq`.
- `RLO=3 RHI=3 GLO=3 GHI=3 BLO=3 BHI=3 NA=1 DEBUG=1 ./csqconv ./small.csq`.
- See `jpegbw` REAME.md to see multiple options that can be passed to this tool.
