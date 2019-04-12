# csqconv

Convert FLIR camera `*.csq` files into a `*.mp4` H.264 video file (you can optionaly save all intermediate files: `*.raw`, `*.jpegls`, `*.png`, `*.hint`, `*.hist` etc.)

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
- Alarms: `RC=1 GC=1 BC=1 RLO=1 RHI=1 GLO=1 GHI=1 BLO=1 BHI=1 NA=1 LIB=libjpegbw.so BF="saturate(x1, .02_1,.98)" GF="saturate(x1, .02, .98)" RF="saturate(x1, .02, .98_1)" csqconv f.csq`
- Alarm and rainbow combined: `RC=1 GC=1 BC=1 RLO=1 RHI=1 GLO=1 GHI=1 BLO=1 BHI=1 NA=1 LIB=libjpegbw.so BF="gsrainbowb(saturate(x1, .01_1,.99))" GF="gsrainbowg(saturate(x1, .01, .99))" RF="gsrainbowr(saturate(x1, 0, .99_1))" csqconv f.csq`.
- Alarm, rainbow and white-hot/black-cold: `RC=1 GC=1 BC=1 RLO=1 RHI=1 GLO=1 GHI=1 BLO=1 BHI=1 NA=1 LIB=libjpegbw.so BF="(3*x1+2*gsrainbowb(saturate(x1, .01_1,.99)))/5" GF="(3*x1+2*gsrainbowg(saturate(x1, .01, .99)))/5" RF="(3*x1+2*gsrainbowr(saturate(x1, 0, .99_1)))/5" csqconv f.csq`.
- Alarms (low B, medium G, high R): `RC=1 GC=1 BC=1 RLO=1 RHI=1 GLO=1 GHI=1 BLO=1 BHI=1 NA=1 LIB=libjpegbw.so RF="saturate(x1, .02, .49)+saturate(x1, .51, .98_1)" GF="saturate(x1, .02, .49_1)-saturate(x1, 0.51,.51_1)+saturate(x1, .51, .98)" BF="saturate(x1, .02_1, .49)+saturate(x1, .51, .98)" csqconv f.csq`.
- See `jpegbw` README.md to see multiple options that can be passed to this tool.

# Other examples

- `RC=1 GC=1 BC=1 RLO=1 RHI=1 GLO=1 GHI=1 BLO=1 BHI=1 NA=1 LIB=libjpegbw.so RF="saturate(((x1-.02)*1.04)^2., .02, .49)+saturate(((x1-.02)*1.04)^2., .51, .98_1)" GF="saturate(((x1-.02)*1.04)^2., .02, .49_1)-saturate(((x1-.02)*1.04)^2., 0.51,.51_1)+saturate(((x1-.02)*1.04)^2., .51, .98)" BF="saturate(((x1-.02)*1.04)^2., .02_1, .49)+saturate(((x1-.02)*1.04)^2., .51, .98)" csqconv FLIR0089.csq`.
- `RC=1 GC=1 BC=1 RLO=1 RHI=1 GLO=1 GHI=1 BLO=1 BHI=1 NA=1 LIB=libjpegbw.so BF="saturate((((1.-x1)-.02)*1.04)^2., .02_1,.98)" GF="saturate((((1-x1)-.02)*1.04)^2., .02, .98)" RF="saturate((((1.-x1)-.02)*1.04)^2., 0, .98_1)" csqconv FLIR0096.csq`.
- `RC=1 GC=1 BC=1 RLO=1 RHI=1 GLO=1 GHI=1 BLO=1 BHI=1 NA=1 LIB=libjpegbw.so RF="(sin(x1*6.283056*5.6+4)+1)/2" GF="(sin(x1*6.283056*4.7+3)+1)/2" BF="(sin(x1*6.283056*3.8+2)+1)/2" csqconv f.csq`.
- `RC=1 GC=1 BC=1 RLO=3 RHI=3 GLO=3 GHI=3 BLO=3 BHI=3 NA=1 LIB=libjpegbw.so RF="1-(cos(x1*3.1415926)+1)/2" GF="x1" BF="x1" csqconv f.csq`.
- `RC=1 GC=1 BC=1 RLO=3 RHI=3 GLO=3 GHI=3 BLO=3 BHI=3 NA=1 LIB=libjpegbw.so RF="(4*x1^2+gsrainbowr(x1))/5" GF="(4*x1^2.5+gsrainbowg(x1))/5" BF="(4*x1^3+gsrainbowb(x1))/5" csqconv f.csq`.
- `RC=1 GC=1 BC=1 RLO=3 RHI=3 GLO=3 GHI=3 BLO=3 BHI=3 NA=1 LIB=libjpegbw.so RF="(sin(2*3.1415926535*(x1-.3))+1)/2" GF="(sin(2*3.1415926535*(x1-.25))+1)/2" BF="(sin(2*3.1415926535*(x1-.2))+1)/2" csqconv f.csq`.
- `RC=1 GC=1 BC=1 RLO=3 RHI=3 GLO=3 GHI=3 BLO=3 BHI=3 NA=1 LIB=libjpegbw.so RF="(sin(1.6*3.1415926535*(x1-.3))+1)/2" GF="(sin(2*3.1415926535*(x1-.25))+1)/2" BF="(sin(2*3.1415926535*(x1-.2))+1)/2" csqconv f.csq`.
- `RC=1 GC=1 BC=1 RLO=1 RHI=1 GLO=1 GHI=1 BLO=1 BHI=1 NA=1 RF='if(x1<.2,5*x1,if(x1<.6,1,1.75-1.25*x1))' GF='if(x1<.2,0,if(x1<.4,5*(x1-.2),if(x1<.8,1,3-2.5*x1)))' BF='if(x1<.4,0,if(x1<.6,5*(x1-.4),1))' csqconv f.csq`.
- `RC=1 GC=1 BC=1 RLO=1 RHI=1 GLO=1 GHI=1 BLO=1 BHI=1 NA=1 RF='if(x1>.5,2*(1-x1),1)' GF='if(x1>.5,0,(.5-x1)*2)' BF='if(x1>.5,0,(.5-x1)*2)' csqconv f.csq`.
- `RC=1 GC=1 BC=1 RLO=3 RHI=3 GLO=3 GHI=3 BLO=3 BHI=3 NA=1 RR=1 RG=0 RB=0 GR=0 GG=1 GB=0 BR=0 BG=0 BB=1 LIB=libjpegbw.so RF="x1" GF="x1" BF="x1" INF=32 EINF=1 csqconv f.csq`.
- `RC=1 GC=1 BC=1 RLO=3 RHI=3 GLO=3 GHI=3 BLO=3 BHI=3 NA=1 LIB=libjpegbw.so RF="(cos(x1*3.1415926)+1)/2" GF="1-x1" BF="1-x1" HPOW=1 INF=20 csqconv f.csq`.
- `RC=1 GC=1 BC=1 RLO=.2 RHI=.2 GLO=.2 GHI=.2 BLO=.2 BHI=.2 NA=1 RR=299 RG=587 RB=114 GR=299 GG=587 GB=114 BR=299 BG=587 BB=114 LIB=libjpegbw.so RF="x1" GF="x1" BF="x1" HPOW=.7 EINF='1' INF=30 csqconv f.csq`.
- `RC=1 GC=1 BC=1 RLO=.2 RHI=.2 GLO=.2 GHI=.2 BLO=.2 BHI=.2 NA=1 LIB=libjpegbw.so RF="((1-x1)*9+gsrainbowr(x1)*1)/10" GF="((1-x1)*9+gsrainbowg(x1)*1)/10" BF="((1-x1)*9+gsrainbowb(x1)*1)/10" INF=16 EINF=1 DEBUG=1 OUTPUT=1 MODE=veryslow CRF=12 csqconv f.csq`.
- Best `RC=1 GC=1 BC=1 RLO=3 RHI=3 GLO=3 GHI=3 BLO=3 BHI=3 NA=1 LIB=libjpegbw.so RF="(cos(x1*3.1415926)+1)/2" GF="1-x1" BF="1-x1" HPOW=1 INF=20 ./csqconv f.csq`.
- Blended rainbow: `RC=1 GC=1 BC=1 RLO=.2 RHI=.2 GLO=.2 GHI=.2 BLO=.2 BHI=.2 NA=1 LIB=libjpegbw.so RF="((1-x1)*9+gsrainbowr(x1)*1)/10" GF="((1-x1)*9+gsrainbowg(x1)*1)/10" BF="((1-x1)*9+gsrainbowb(x1)*1)/10" INF=16 csqconv f.csq`.
- Info rainbow: `RC=1 GC=1 BC=1 RLO=3 RHI=3 GLO=3 GHI=3 BLO=3 BHI=3 NA=1 LIB=libjpegbw.so RF="(x1*6+gsrainbowr(x1)*1)/7" GF="(x1*6+gsrainbowg(x1)*1)/7" BF="(x1*6+gsrainbowb(x1)*1)/7" INF=16 csqconv f.csq`
- Use `HINT=1` to calculate intensity range for 32 frames (less bumpy exposure adjustments), add `MF=60` to average instensity from 60 frames. Use `MODE=veryslow CRF=1` to have best possible loosy compression.
- Use `MODE=veryslow CRF=0` to have X264 loseless compression (huge files), use `MODE=mpng` to have totally loseless montion PNG compression. Use `PQ=n` to control PNG compression levels: 0=default, 1=no compression, 2=best speed, 3=best size.
- Use `OGS=1` to enable grayscale PNG outputs (so video will be made from 16-bit PNGs instead of 64-bit RGBA PNGs). Use `GSR=2 GSG=7 GSB=1` to set R,G,B channels mix ratio for grayscale `OGS` output.
- So to have a really loseless video but with smallest possible size use: `INF=16 HINT=1 HINTREQ=1 MF=46 MODE=mpng PQ=3 OGS=1 GSR=0.2125 GSG=0.7154 GSB=0.0721 csqconv f.csq`. If you want loseless but not using mpng replace `MODE=mpng` with `MODE=veryslow CRF=0`.
- If you want best possible quality but not loseless replace `MODE=mpng` with `MODE=veryslow CRF=1`. `CRF` values up to 12 or even 15 will give visually the same quality but a lot smaller sizes.
- Setting `PQ=3` option only makes sens with using `MODE=mpng` otherwise the X264 encoder will take care of all intermediate frames encoding, so it is even suggested to use `PQ=1` which means no compression and maximal speed.
- Typical REAL loseless: `RC=1 GC=1 BC=1 RLO=1 RHI=1 GLO=1 GHI=1 BLO=1 BHI=1 NA=1 INF=16 HINT=1 PQ=3 OGS=1 RF="1-x1" GF="1-x1" BF="1-x1" MODE=veryslow CRF=0 csqconv f.csq`.
- Hot and Cold black - medium white, grayscaled using X.264 and skipping intermediate PNG compression: `LIB=libjpegbw.so RC=1 GC=1 BC=1 RLO=1 RHI=1 GLO=1 GHI=1 BLO=1 BHI=1 NA=1 INF=18 HINT=1 PQ=1 OGS=1 RF="sin(3.1415926535*x1)" GF="sin(3.1415926535*x1)" BF="sin(3.1415926535*x1)" MODE=veryslow csqconv f.csq`.
- Double blended rainbow: `RC=1 GC=1 BC=1 RLO=3 RHI=3 GLO=3 GHI=3 BLO=3 BHI=3 NA=1 LIB=libjpegbw.so RF="(if(2*x1>1,2*x1-1,2*x1)*6+gsrainbowr(if(2*x1>1,2*x1-1,2*x1))*1)/7" GF="(if(2*x1>1,2*x1-1,2*x1)*6+gsrainbowg(if(2*x1>1,2*x1-1,2*x1))*1)/7" BF="(if(2*x1>1,2*x1-1,2*x1)*6+gsrainbowb(if(2*x1>1,2*x1-1,2*x1))*1)/7" INF=16 HINT=1 csqconv f.csq`.
- Rainblow with gray too: `LIB=libjpegbw.so RC=1 GC=1 BC=1 RLO=.1 RHI=.1 GLO=.1 GHI=.1 BLO=.1 BHI=.1 NA=1 INF=20 HINT=1 PQ=1 MODE=veryslow RF='if(x1<.1,10*(.1-x1),if(x1>.9,(1.-x1)*10,gsrainbowr((x1-.1)*1.25)))' GF='if(x1<.1,10*(.1-x1),if(x1>.9,(1.-x1)*10,gsrainbowg((x1-.1)*1.25)))' BF='if(x1<.1,10*(.1-x1),if(x1>.9,(1.-x1)*10,gsrainbowb((x1-.1)*1.25)))' csqconv f.csq`.
- Fixed from ~10C to ~40C with alarms: `LIB=libjpegbw.so RC=1 GC=1 BC=1 RLOI=14000 RHII=23400 GLOI=14000 GHII=23400 BLOI=14000 BHII=23400 NA=1 INF=16 PQ=1 RF="if(x1>.99,1,if(x1<.01,0,x1))" GF="if(x1>.99,0,if(x1<.01,0,x1))" BF="if(x1>.99,0,if(x1<.01,1,x1))" csqconv f.csq`.
- About 3.8C to 20.3C with alarms: `LIB=libjpegbw.so RC=1 GC=1 BC=1 RLOI=13000 RHII=22000 GLOI=13000 GHII=22000 BLOI=13000 BHII=22000 NA=1 INF=16 PQ=1 RF="if(x1>.99,1,if(x1<.01,0,x1))" GF="if(x1>.99,0,if(x1<.01,0,x1))" BF="if(x1>.99,0,if(x1<.01,1,x1))" csqconv f.csq`

# My guess work

This is my guess work between converting Unsigned 16bit int values (u16) into temperatures in Celsius. camera has two modes: `-20 - 120` (but it shows temperatures down to -40 so I've used two guesses) and `0 - 650` - I've marked it with `(*)`.
- Use `ttou16.sh temp_in_celsius_1 temp_in_celsius_2 ...` to get UINT16 values.
- Use `u16tot.sh U16_value_1 U16_value_2 ...` to get temperatures in Celsius.
