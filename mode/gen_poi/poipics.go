package poi

import "image"

type GetPoiPic func() image.NRGBA

var PixelCount = 58

var PoiPicsCount = 6

var PoiPicsGetter = []GetPoiPic{GetPoiPic_fire, GetPoiPic_geo, GetPoiPic_honeycorb, GetPoiPic_man_flower, GetPoiPic_rainbow, GetPoiPic_wave, }